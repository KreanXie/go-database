package internal

import (
	"os"
	"testing"
)

func setupDiskManager(t *testing.T) *DiskManager {
	// 创建一个临时文件来模拟磁盘文件
	dbFileName := "test_db.db"
	logFileName := "test_log.db"

	dbFile, err := os.Create(dbFileName)
	if err != nil {
		t.Fatalf("无法创建测试数据库文件: %v", err)
	}
	dbFile.Truncate(4096 * 10)

	logFile, err := os.Create(logFileName)
	if err != nil {
		t.Fatalf("无法创建测试日志文件: %v", err)
	}
	logFile.Truncate(4096 * 10)

	dm := &DiskManager{
		DBFile:      dbFile,
		LogFile:     logFile,
		DBFileName:  dbFileName,
		LogFileName: logFileName,
	}
	return dm
}

func cleanupDiskManager(dm *DiskManager) {
	dm.DBFile.Close()
	dm.LogFile.Close()
	if err := os.Remove(dm.DBFileName); err != nil {
		panic(err)
	}
	if err := os.Remove(dm.LogFileName); err != nil {
		panic(err)
	}
	dm.ShutDown()
}

func TestBufferPool(t *testing.T) {
	dm := setupDiskManager(t)
	defer cleanupDiskManager(dm)

	poolSize := 3
	pageSize := 4096
	k := 2
	bm := NewBufferPoolManager(dm, poolSize, pageSize, k)

	// 测试 FetchPage
	pageID := 1
	p, err := bm.FetchPage(pageID)
	if err != nil {
		t.Fatalf("FetchPage 失败: %v", err)
	}
	if p.PageID != pageID {
		t.Fatalf("FetchPage 返回的 PageID 不匹配: 期望 %d, 但得到 %d", pageID, p.PageID)
	}

	// 测试 UnpinPage
	err = bm.UnpinPage(pageID, true)
	if err != nil {
		t.Fatalf("UnpinPage 失败: %v", err)
	}

	// 测试 FlushPage
	err = bm.FlushPage(pageID)
	if err != nil {
		t.Fatalf("FlushPage 失败: %v", err)
	}

	// 测试 NewPage
	newPageID, err := bm.NewPage()
	if err != nil {
		t.Fatalf("NewPage 失败: %v", err)
	}
	if newPageID <= 0 {
		t.Fatalf("NewPage 返回的 ID 无效: %d", newPageID)
	}

	// 测试 DeletePage
	err = bm.DeletePage(newPageID)
	if err != nil {
		t.Fatalf("DeletePage 失败: %v", err)
	}
}
