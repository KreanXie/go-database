package internal

import "sync"

const DefaultPageSize = PageSize

// Page 页面结构
type Page struct {
	PageID   int
	Data     []byte
	IsDirty  bool
	PinCount int
	mu       sync.Mutex
}
