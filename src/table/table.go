package table

import (
	"sync"

	bufferpool "go-database/src/buffer_pool"
	"go-database/src/disk"
	"go-database/src/row"
)

type Table struct {
	Name          string
	Schema        *Schema
	BufferPoolMgr *bufferpool.Manager
	DiskMgr       *disk.Manager
	mu            sync.Mutex
}

func NewTable(name string, schema *Schema, bufferPoolMgr *bufferpool.Manager, diskMgr *disk.Manager) *Table {
	return &Table{
		Name:          name,
		Schema:        schema,
		BufferPoolMgr: bufferPoolMgr,
		DiskMgr:       diskMgr,
	}
}

func (t *Table) Insert(row *row.Row) error {
	// 1. 查找可用的 Page
	pageID := t.findAvailablePage() // 这里需要维护一个数据结构来存储可用的 Page
	if pageID == -1 {
		// 2. 若无可用 Page，申请新 Page
		newPageID, err := t.BufferPoolMgr.NewPage()
		if err != nil {
			return err
		}
		pageID = newPageID
	}

	// 3. 获取 Page
	page, err := t.BufferPoolMgr.FetchPage(pageID)
	if err != nil {
		return err
	}
	defer t.BufferPoolMgr.UnpinPage(pageID, true) // 释放 Page，并标记为脏页

	// 4. 插入数据
	if _, err := page.InsertRow(row); err != nil {
		return err
	}

	return nil
}

func (t *Table) findAvailablePage() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	for pageID, page := range t.BufferPoolMgr.PageTable {
		// 获取 Page
		page, err := t.BufferPoolMgr.FetchPage(pageID)
		if err != nil {
			continue
		}

		// 检查 Page 是否有可用空间
		if page.HasFreeSpace() {
			t.BufferPoolMgr.UnpinPage(pageID, false) // 释放 Page（未修改）
			return pageID
		}

		// 释放 Page（未修改）
		t.BufferPoolMgr.UnpinPage(pageID, false)
	}
	return -1 // 没有可用 Page
}
