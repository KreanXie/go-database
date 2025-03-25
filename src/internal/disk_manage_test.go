package internal

import (
	"log"
	"os"
	"testing"
)

const (
	testPageSize = PageSize
)

func TestMain(m *testing.M) {
	if _, err := os.Create(FilePath + "test.db"); err != nil {
		panic(err)
	}

	if _, err := os.Create(FilePath + "test.log"); err != nil {
		panic(err)
	}

	defer cleanup()

	m.Run()
}

func TestWritePage(t *testing.T) {
	dm, err := NewDiskManager("test.db")
	if err != nil {
		t.Fatalf("Failed to create DiskManager: %v", err)
	}

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
	dm, err := NewDiskManager("test.db")
	if err != nil {
		t.Fatalf("Failed to create DiskManager: %v", err)
	}
	defer cleanup()

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

// cleanup 清理测试文件
func cleanup() {
	if err := os.Remove(FilePath + "test.db"); err != nil {
		log.Printf("Failed to remove database file: %v", err)
	}

	if err := os.Remove(FilePath + "test.log"); err != nil {
		log.Printf("Failed to remove log file: %v", err)
	}
}
