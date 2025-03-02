package disk

import (
	"errors"
	"os"
	"sync"
)

const (
	// PageSize 页面大小
	PageSize = 4096

	// FilePath 文件路径
	FilePath = "../data/"
)

type DiskManager struct {
	dbFile       *os.File
	logFile      *os.File
	dbFileName   string
	logFileName  string
	mu           sync.Mutex
	numWrites    int
	pageCapacity int
}

// NewDiskManager 构造函数
func NewDiskManager(dbFileName string) (*DiskManager, error) {
	dbFile, err := os.OpenFile(FilePath+dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.New("failed to open database file")
	}

	logFileName := dbFileName + ".log"
	logFile, err := os.OpenFile(FilePath+logFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.New("failed to open log file")
	}

	return &DiskManager{
		dbFile:      dbFile,
		logFile:     logFile,
		dbFileName:  dbFileName,
		logFileName: logFileName,
		mu:          sync.Mutex{},
	}, nil
}

// WritePage 将数据写入文件
func (dm *DiskManager) WritePage(pageID int, pageData []byte) error {
	if len(pageData) != PageSize {
		return errors.New("page data size dismatch")
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	offset := int64(pageID) * PageSize
	_, err := dm.dbFile.WriteAt(pageData, offset)
	if err != nil {
		return errors.New("failed to write page data: " + err.Error())
	}

	dm.numWrites++
	return nil
}

// ReadPage 读取页
func (dm *DiskManager) ReadPage(pageID int, pageData []byte) error {
	if len(pageData) != PageSize {
		return errors.New("page data size dismatch")
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	offset := int64(pageID) * PageSize
	_, err := dm.dbFile.ReadAt(pageData, offset)
	if err != nil {
		return errors.New("failed to read db file: " + err.Error())
	}

	return nil
}

// ShutDown 关闭文件流
func (dm *DiskManager) ShutDown() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if err := dm.dbFile.Close(); err != nil {
		return errors.New("failed to close db file: " + err.Error())
	}
	if err := dm.logFile.Close(); err != nil {
		return errors.New("failed to close log file: " + err.Error())
	}

	return nil
}
