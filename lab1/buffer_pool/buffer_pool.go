package bufferpool

import (
	"container/list"

	"cmu-database/lab1/disk"
	lruk "cmu-database/lab1/lru_k"
)

const PageSize = disk.PageSize

// Page 页面结构
type Page struct {
	PageID   int
	Data     []byte
	IsDirty  bool
	PinCount int
}

// BufferPoolManager 缓冲池管理器
type BufferPoolManager struct {
	diskManager *disk.DiskManager
	poolSize    int
	pages       map[int]*list.Element
	lruList     lruk.LRUKReplacer
}

func NewBufferPoolManager() *BufferPoolManager {

}
