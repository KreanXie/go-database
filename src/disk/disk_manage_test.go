package disk

import (
	"os"
	"testing"
)

const (
	testDBFileName = "test.db"
	testPageSize   = PageSize
)

func TestMain(m *testing.M) {
	err := os.Mkdir("../data", 0755)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll("../data")

	m.Run()
}

func TestNewManager(t *testing.T) {
	// 初始化测试
	dm, err := NewManager(testDBFileName)
	if err != nil {
		t.Fatalf("Failed to create Manager: %v", err)
	}
	defer cleanup(dm, t)

	// 验证数据库文件是否创建
	if _, err := os.Stat(FilePath + testDBFileName); os.IsNotExist(err) {
		t.Fatalf("Database file not created: %v", err)
	}

	// 验证日志文件是否创建
	if _, err := os.Stat(FilePath + testDBFileName + ".log"); os.IsNotExist(err) {
		t.Fatalf("Log file not created: %v", err)
	}
}

func TestWritePage(t *testing.T) {
	dm, err := NewManager(testDBFileName)
	if err != nil {
		t.Fatalf("Failed to create Manager: %v", err)
	}
	defer cleanup(dm, t)

	pageID := 0
	pageData := make([]byte, PageSize)
	for i := range pageData {
		pageData[i] = byte(i % 256)
	}

	// 写入页面数据
	err = dm.WritePage(pageID, pageData)
	if err != nil {
		t.Fatalf("Failed to write page: %v", err)
	}

	// 验证页面是否写入正确
	readData := make([]byte, PageSize)
	err = dm.ReadPage(pageID, readData)
	if err != nil {
		t.Fatalf("Failed to read page after writing: %v", err)
	}

	if string(pageData) != string(readData) {
		t.Fatalf("Page data mismatch after write and read")
	}
}

func TestReadPage(t *testing.T) {
	dm, err := NewManager(testDBFileName)
	if err != nil {
		t.Fatalf("Failed to create Manager: %v", err)
	}
	defer cleanup(dm, t)

	pageID := 1
	pageData := make([]byte, PageSize)
	for i := range pageData {
		pageData[i] = byte((i + 1) % 256)
	}

	// 写入页面数据
	err = dm.WritePage(pageID, pageData)
	if err != nil {
		t.Fatalf("Failed to write page: %v", err)
	}

	// 读取页面数据
	readData := make([]byte, PageSize)
	err = dm.ReadPage(pageID, readData)
	if err != nil {
		t.Fatalf("Failed to read page: %v", err)
	}

	if string(pageData) != string(readData) {
		t.Fatalf("Page data mismatch")
	}
}

func TestShutDown(t *testing.T) {
	dm, err := NewManager(testDBFileName)
	if err != nil {
		t.Fatalf("Failed to create Manager: %v", err)
	}

	// 关闭
	cleanup(dm, t)

	// 尝试写入页面，验证关闭状态
	err = dm.WritePage(0, make([]byte, PageSize))
	if err == nil {
		t.Fatalf("WritePage should fail after ShutDown")
	}
}

// cleanup 清理测试文件
func cleanup(dm *Manager, t *testing.T) {
	err := dm.ShutDown()
	if err != nil {
		t.Logf("Failed to shut down Manager: %v", err)
	}
	if err := os.Remove(FilePath + testDBFileName); err != nil {
		t.Logf("Failed to remove database file: %v", err)
	}
	if err := os.Remove(FilePath + testDBFileName + ".log"); err != nil {
		t.Logf("Failed to remove log file: %v", err)
	}
}
