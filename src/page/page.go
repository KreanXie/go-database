package page

import (
	"go-database/src/disk"
)

const DefaultPageSize = disk.PageSize

// Page 页面结构
type Page struct {
	PageID   int
	Data     []byte
	IsDirty  bool
	PinCount int
}
