package fileCache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"
)




func TestCache(t *testing.T) {
	var testKey = "testKey"
	var testValue = time.Now().Format(time.RFC3339)
	cache,err := Init(".","test")
	if err!=nil {
		t.Error(err)
	}
	//测试创建
	cache.Put(testKey,testValue,1000*5)
	if v,ok:=cache.Get(testKey);v!=testValue || !ok {
		t.Error("cache put error")
	}

	//测试删除
	cache.Del(testKey)
	if v,ok:=cache.Get(testKey);v==testValue || ok {
		t.Error("cache del error")
	}

	//测试超时
	cache.Put(testKey,testValue,500)
	time.Sleep(time.Second)
	if v,ok:=cache.Get(testKey);v==testValue || ok {
		t.Error("cache timeout error")
	}
	//测试将缓存写入文件
	cache.Put(testKey,testValue,5*60*1000)
	time.Sleep(2*time.Second)
	tempByte,err := ioutil.ReadFile(cache.DbFilePath)
	var readData map[string]string
	json.Unmarshal(tempByte,&readData)
	if readData[testKey]!=testValue {
		t.Error("cache save file error")
	}
	//测试从文件读取缓存
	cache.DbData= map[string]string{}
	file2Map(cache)
	if value,ok:=cache.Get(testKey);!ok||value!=testValue {
		t.Error("load file error")
	}
	os.Remove(cache.DbFilePath)
	os.Remove(cache.ExpireFilePath)
}