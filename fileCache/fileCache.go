package fileCache

import (
	"bytes"
	"encoding/gob"
	"github.com/ayu-666/gopkg/dir"
	"github.com/ayu-666/gopkg/timeTask"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	dbFileSuffix     = ".db.data"
	expireFileSuffix = ".expire.data"
)

type FCInterface interface {
	Put(key string, value string, millisecondTTL int64, dbName string) error
	Get(key string, dbName string) (string, error)
	Del(key string, dbName string) error
}

type FcDb struct {
	DbRwLock       sync.RWMutex
	DbData         map[string]string
	ExpireData     map[string]int64
	DbName         string
	DbFilePath     string
	ExpireFilePath string
	fileLock       sync.Mutex
	taskChannel    chan int
}

func Init(dirPath string, dbName string) (*FcDb, error) {
	var err error
	if dirPath == "" {
		dirPath = "db"
	}
	err = dir.Mkdir(dirPath, 0766)
	if err != nil {
		return nil, err
	}
	db := &FcDb{
		DbName:         dbName,
		DbData:         map[string]string{},
		ExpireData:     map[string]int64{},
		DbRwLock:       sync.RWMutex{},
		fileLock:       sync.Mutex{},
		taskChannel:    make(chan int, 1),
		DbFilePath:     filepath.Join(dirPath, dbName+dbFileSuffix),
		ExpireFilePath: filepath.Join(dirPath, dbName+expireFileSuffix),
	}
	err = file2Map(db)
	if err != nil {
		return nil, err
	}
	timeTask.NewTask().SetInterval(100 * time.Millisecond).SetConsumer(func() { //设置消费者
		Map2FileConsumer(db)
	})
	return db, err
}

func file2Map(db *FcDb) (err error) {
	db.DbRwLock.Lock()
	db.fileLock.Lock()
	defer func() {
		db.DbRwLock.Unlock()
		db.fileLock.Unlock()
	}()
	var wg sync.WaitGroup
	dbData := map[string]string{}
	expireData := map[string]int64{}
	var expireErr error
	var dbErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		dbByte := []byte("")
		dbByte, _ = ioutil.ReadFile(db.DbFilePath)
		rd := bytes.NewReader(dbByte)
		dr := gob.NewDecoder(rd)
		dr.Decode(&dbData)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ExpireByte := []byte("")
		ExpireByte, _ = ioutil.ReadFile(db.ExpireFilePath)
		rd := bytes.NewReader(ExpireByte)
		dr := gob.NewDecoder(rd)
		dr.Decode(&expireData)
	}()

	wg.Wait()
	if expireErr != nil {
		err = expireErr
		return
	}
	if dbErr != nil {
		err = dbErr
		return
	}
	db.DbData = dbData
	db.ExpireData = expireData
	return
}

func Map2File(db *FcDb) (err error) {
	db.fileLock.Lock()
	db.DbRwLock.RLock()
	defer func() {
		db.fileLock.Unlock()
		db.DbRwLock.RUnlock()
	}()
	var dbErr error
	var expireErr error
	//读取数据文件
	_dbFile, err := os.OpenFile(db.DbFilePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0766)
	defer _dbFile.Close()

	if err != nil {
		return
	}
	//读取过期列表文件
	_expireFile, err := os.OpenFile(db.ExpireFilePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0766)
	defer _expireFile.Close()

	if err != nil {
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		ed := gob.NewEncoder(_dbFile)
		dbErr = ed.Encode(db.DbData)
	}()
	go func() {
		defer wg.Done()
		ed := gob.NewEncoder(_expireFile)
		expireErr = ed.Encode(db.ExpireData)
	}()
	wg.Wait()
	if expireErr != nil {
		err = expireErr
	}
	if dbErr != nil {
		err = dbErr
	}
	return
}

func (_this *FcDb) Get(key string) (value string, exist bool) {
	_this.DbRwLock.RLock()
	defer func() {
		_this.DbRwLock.RUnlock()
	}()
	value = ""
	exist = false
	expireTime := _this.ExpireData[key]
	currTime := time.Now().UnixMilli()
	if currTime > expireTime {
		defer _this.DbRwLock.Unlock()
		_this.DbRwLock.Lock()
		delete(_this.ExpireData, key)
		delete(_this.DbData, key)
		Map2FileProducer(_this)
		return
	}
	value, exist = _this.DbData[key]
	return
}

func (_this *FcDb) Del(key string) {
	_this.DbRwLock.Lock()
	defer func() {
		_this.DbRwLock.Unlock()
		Map2FileProducer(_this)
	}()
	delete(_this.ExpireData, key)
	delete(_this.DbData, key)
}

func (_this *FcDb) Put(key string, value string, milliSecond int64) {
	_this.DbRwLock.Lock()
	defer func() {
		_this.DbRwLock.Unlock()
		Map2FileProducer(_this)
	}()
	_this.DbData[key] = value
	_this.ExpireData[key] = time.Now().UnixMilli() + milliSecond
}

func Map2FileConsumer(db *FcDb) {
	select {
	case <-db.taskChannel:
		Map2File(db)
	default:
	}
}
func Map2FileProducer(db *FcDb) {
	select {
	case db.taskChannel <- 1:
	default:
	}
}
