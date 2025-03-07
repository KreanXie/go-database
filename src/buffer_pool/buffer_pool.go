package bufferpool

import (
	"fmt"
	"sync"

	"go-database/src/disk"
	"go-database/src/lruk"
	"go-database/src/page"
)

// Manager 缓冲池管理器
type Manager struct {
	DiskManager *disk.Manager
	Replacer    *lruk.Replacer
	PageTable   map[int]*page.Page
	PinnedPages map[int]int
	PoolSize    int
	PageSize    int
	mu          sync.Mutex
}

// NewManager 创建一个新的 Manager 实例
func NewManager(diskManager *disk.Manager, poolSize, DefaultPageSize, k int) *Manager {
	replacer := lruk.NewReplacer(poolSize, k)
	return &Manager{
		DiskManager: diskManager,
		Replacer:    replacer,
		PageTable:   make(map[int]*page.Page),
		PinnedPages: make(map[int]int),
		PoolSize:    poolSize,
		PageSize:    DefaultPageSize,
	}
}

// FetchPage 从缓冲池或磁盘中获取指定页面
func (m *Manager) FetchPage(pageID int) (*page.Page, error) {
	// 尝试从缓冲池获取page
	if p, ok := m.PageTable[pageID]; ok {
		return p, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果池子的大小超过的预定大小，则驱逐一个page
	if len(m.PageTable) >= m.PoolSize {
		if err := m.evictPage(); err != nil {
			return nil, err
		}
	}

	// 创建一个新page加入到缓冲池
	p := &page.Page{
		PageID: pageID,
		Data:   make([]byte, page.DefaultPageSize),
	}

	// 从磁盘中读取数据
	if err := m.DiskManager.ReadPage(int(pageID), p.Data); err != nil {
		return nil, err
	}

	// 加入到缓冲池
	m.PageTable[pageID] = p
	p.PinCount++

	// 在lru replacer中记录这个page
	if err := m.Replacer.RecordAccess(int(pageID), 0); err != nil {
		return nil, err
	}

	return p, nil
}

// UnpinPage 解除固定页面并处理相关的脏页面写回
// 通常在FetchPage，并结束对页的操作之后，需要UnpinPage
func (m *Manager) UnpinPage(pageID int, isDirty bool) error {
	// 确认该页面是否存在
	p, ok := m.PageTable[pageID]
	if !ok {
		return fmt.Errorf("page %d not found", pageID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 页面没有被标记
	if p.PinCount == 0 {
		return fmt.Errorf("page %d is already unpinned", pageID)
	}

	// 取消标记
	p.PinCount--

	// 页面是脏页，需要写回磁盘
	if isDirty {
		p.IsDirty = true
	}

	// 如果这个页没有被标记且是脏的，则写回磁盘
	if p.PinCount == 0 && p.IsDirty {
		if err := m.FlushPage(pageID); err != nil {
			return err
		}
	}

	return nil
}

// FlushPage 将指定页面的数据刷新到磁盘
func (m *Manager) FlushPage(pageID int) error {
	// 确认该页面是否存在
	p, ok := m.PageTable[pageID]
	if !ok {
		return fmt.Errorf("page %d not found", pageID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果是脏页，则将其写回磁盘
	if p.IsDirty {
		if err := m.DiskManager.WritePage(int(pageID), p.Data); err != nil {
			return err
		}
		// 重置为干净
		p.IsDirty = false
	}

	return nil
}

// NewPage 创建一个新页到缓冲池，并放回id
func (m *Manager) NewPage() (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pageID := int(len(m.PageTable) + 1)

	p := &page.Page{
		PageID: pageID,
		Data:   make([]byte, page.DefaultPageSize),
	}

	m.PageTable[pageID] = p

	return pageID, nil
}

// DeletePage 从缓冲池中删除页
func (m *Manager) DeletePage(pageID int) error {
	// 确认该页是否存在
	p, ok := m.PageTable[pageID]
	if !ok {
		return fmt.Errorf("page %d not found", pageID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果页面被标记，表示有进程正在使用这个页，不予删除
	if p.PinCount > 0 {
		return fmt.Errorf("page %d is already pinned", pageID)
	}

	// 如果是脏页，则先刷新到磁盘，再删除
	if p.IsDirty {
		if err := m.FlushPage(pageID); err != nil {
			return err
		}
	}

	delete(m.PageTable, pageID)

	return nil
}

// FlushAllPages 刷新所有的脏页到磁盘
func (m *Manager) FlushAllPages() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for pageID, p := range m.PageTable {
		if p.IsDirty {
			if err := m.FlushPage(pageID); err != nil {
				return err
			}
		}
	}

	return nil
}

// evictPage 驱逐一个页
func (m *Manager) evictPage() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	evictedFrameID, err := m.Replacer.Evict(-1)
	if err != nil {
		return err
	}

	for pageID, p := range m.PageTable {
		if int(p.PageID) == evictedFrameID {
			if p.IsDirty {
				if err := m.FlushPage(pageID); err != nil {
					return err
				}
			}

			delete(m.PageTable, pageID)
			break
		}
	}

	return nil
}
