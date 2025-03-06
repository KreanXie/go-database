package disk

import (
	"fmt"
	"os"
	"sync"
)

const (
	// PageSize 页面大小
	PageSize = 4096

	// FilePath 文件路径
	FilePath = "../data/"
)

type Manager struct {
	dbFile       *os.File
	logFile      *os.File
	dbFileName   string
	logFileName  string
	mu           sync.Mutex
	numWrites    int
	pageCapacity int
}

// NewManager 构造函数
func NewManager(dbFileName string) (*Manager, error) {
	dbFile, err := os.OpenFile(FilePath+dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open db file: %v", err)
	}

	logFileName := dbFileName + ".log"
	logFile, err := os.OpenFile(FilePath+logFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &Manager{
		dbFile:      dbFile,
		logFile:     logFile,
		dbFileName:  dbFileName,
		logFileName: logFileName,
		mu:          sync.Mutex{},
	}, nil
}

// WritePage 将数据写入文件
func (dm *Manager) WritePage(pageID int, pageData []byte) error {
	if len(pageData) != PageSize {
		return fmt.Errorf("invalid page size")
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	offset := int64(pageID) * PageSize
	_, err := dm.dbFile.WriteAt(pageData, offset)
	if err != nil {
		return fmt.Errorf("write page error: %v", err)
	}

	dm.numWrites++
	return nil
}

// ReadPage 读取页
func (dm *Manager) ReadPage(pageID int, pageData []byte) error {
	if len(pageData) != PageSize {
		return fmt.Errorf("invalid page size")
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	offset := int64(pageID) * PageSize
	n, err := dm.dbFile.ReadAt(pageData, offset)
	if err != nil {
		return fmt.Errorf("read page error: %v", err)
	}

	if n < PageSize {
		return fmt.Errorf("incomplete page read: expected %d bytes, got %d", PageSize, n)
	}

	return nil
}

// ShutDown 关闭文件流
func (dm *Manager) ShutDown() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if err := dm.dbFile.Close(); err != nil {
		return fmt.Errorf("failed to close db file: %v", err)
	}

	if err := dm.logFile.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %v", err)
	}

	return nil
}
