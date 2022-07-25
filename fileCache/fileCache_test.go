package fileCache

import (
	"os"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	var testKey = "testKey"
	var testValue = time.Now().Format(time.RFC3339)
	cache, err := Init(".", "test")
	if err != nil {
		t.Error(err)
	}
	//测试创建
	cache.Put(testKey, testValue, 1000*5)
	if v, ok := cache.Get(testKey); v != testValue || !ok {
		t.Error("cache put error")
	}

	//测试删除
	cache.Del(testKey)
	if v, ok := cache.Get(testKey); v == testValue || ok {
		t.Error("cache del error")
	}

	//测试超时
	cache.Put(testKey, testValue, 500)
	time.Sleep(time.Second)
	if v, ok := cache.Get(testKey); v == testValue || ok {
		t.Error("cache timeout error")
	}
	cache.Put(testKey, testValue, 5*60*1000)
	err = Map2File(cache)
	if err != nil {
		t.Error("save file failed")
	}
	time.Sleep(time.Second)
	err = file2Map(cache)
	if err != nil {
		t.Error("read file failed")
	}
	time.Sleep(time.Second)
	os.Remove(cache.DbFilePath)
	os.Remove(cache.ExpireFilePath)
}
