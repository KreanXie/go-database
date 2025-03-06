package bufferpool

import (
	"fmt"
	"sync"

	"go-database/src/disk"
	"go-database/src/lruk"
)

const PageSize = disk.PageSize

// Page 页面结构
type Page struct {
	pageID   int
	data     []byte
	isDirty  bool
	pinCount int
}

// Manager 缓冲池管理器
type Manager struct {
	diskManager *disk.Manager
	replacer    *lruk.Replacer
	pageTable   map[int]*Page
	pinnedPages map[int]int
	poolSize    int
	pageSize    int
	mu          sync.Mutex
}

// NewManager 创建一个新的 Manager 实例
func NewManager(diskManager *disk.Manager, poolSize, pageSize, k int) *Manager {
	replacer := lruk.NewReplacer(poolSize, k)
	return &Manager{
		diskManager: diskManager,
		replacer:    replacer,
		pageTable:   make(map[int]*Page),
		pinnedPages: make(map[int]int),
		poolSize:    poolSize,
		pageSize:    pageSize,
	}
}

// FetchPage 从缓冲池或磁盘中获取指定页面
func (m *Manager) FetchPage(pageID int) (*Page, error) {
	// 尝试从缓冲池获取page
	if page, ok := m.pageTable[pageID]; ok {
		return page, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果池子的大小超过的预定大小，则驱逐一个page
	if len(m.pageTable) >= m.poolSize {
		if err := m.evictPage(); err != nil {
			return nil, err
		}
	}

	// 创建一个新page加入到缓冲池
	page := &Page{
		pageID: pageID,
		data:   make([]byte, PageSize),
	}

	// 从磁盘中读取数据
	if err := m.diskManager.ReadPage(int(pageID), page.data); err != nil {
		return nil, err
	}

	// 加入到缓冲池
	m.pageTable[pageID] = page
	page.pinCount++

	// 在lru replacer中记录这个page
	if err := m.replacer.RecordAccess(int(pageID), 0); err != nil {
		return nil, err
	}

	return page, nil
}

// UnpinPage 解除固定页面并处理相关的脏页面写回
// 通常在FetchPage，并结束对页的操作之后，需要UnpinPage
func (m *Manager) UnpinPage(pageID int, isDirty bool) error {
	// 确认该页面是否存在
	page, ok := m.pageTable[pageID]
	if !ok {
		return fmt.Errorf("page %d not found", pageID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 页面没有被标记
	if page.pinCount == 0 {
		return fmt.Errorf("page %d is already unpinned", pageID)
	}

	// 取消标记
	page.pinCount--

	// 页面是脏页，需要写回磁盘
	if isDirty {
		page.isDirty = true
	}

	// 如果这个页没有被标记且是脏的，则写回磁盘
	if page.pinCount == 0 && page.isDirty {
		if err := m.FlushPage(pageID); err != nil {
			return err
		}
	}

	return nil
}

// FlushPage 将指定页面的数据刷新到磁盘
func (m *Manager) FlushPage(pageID int) error {
	// 确认该页面是否存在
	page, ok := m.pageTable[pageID]
	if !ok {
		return fmt.Errorf("page %d not found", pageID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果是脏页，则将其写回磁盘
	if page.isDirty {
		if err := m.diskManager.WritePage(int(pageID), page.data); err != nil {
			return err
		}
		// 重置为干净
		page.isDirty = false
	}

	return nil
}

// NewPage 创建一个新页到缓冲池，并放回id
func (m *Manager) NewPage() (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pageID := int(len(m.pageTable) + 1)

	page := &Page{
		pageID: pageID,
		data:   make([]byte, PageSize),
	}

	m.pageTable[pageID] = page

	return pageID, nil
}

// DeletePage 从缓冲池中删除页
func (m *Manager) DeletePage(pageID int) error {
	// 确认该页是否存在
	page, ok := m.pageTable[pageID]
	if !ok {
		return fmt.Errorf("page %d not found", pageID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果页面被标记，表示有进程正在使用这个页，不予删除
	if page.pinCount > 0 {
		return fmt.Errorf("page %d is already pinned", pageID)
	}

	// 如果是脏页，则先刷新到磁盘，再删除
	if page.isDirty {
		if err := m.FlushPage(pageID); err != nil {
			return err
		}
	}

	delete(m.pageTable, pageID)

	return nil
}

// FlushAllPages 刷新所有的脏页到磁盘
func (m *Manager) FlushAllPages() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for pageID, page := range m.pageTable {
		if page.isDirty {
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

	evictedFrameID, err := m.replacer.Evict(-1)
	if err != nil {
		return err
	}

	for pageID, page := range m.pageTable {
		if int(page.pageID) == evictedFrameID {
			if page.isDirty {
				if err := m.FlushPage(pageID); err != nil {
					return err
				}
			}

			delete(m.pageTable, pageID)
			break
		}
	}

	return nil
}
