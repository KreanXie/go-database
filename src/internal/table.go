package internal

import (
	"sync"
)

type Table struct {
	Name          string
	Schema        *Schema
	BufferPoolMgr *BufferPoolManager
	DiskMgr       *DiskManager
	mu            sync.Mutex
}
