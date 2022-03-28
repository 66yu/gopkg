package weCom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func Post(url string, jsonParams []byte) (map[string]interface{}, error) {
	log.Printf("%s", jsonParams)
	var r map[string]interface{}
	reader := bytes.NewReader(jsonParams)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Println(err.Error())
		return r, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return r, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return requestHandler(&respBytes)
}

func Get(url string) (map[string]interface{}, error) {
	var r map[string]interface{}
	resp, err := http.Get(url)
	if err != nil {
		return r, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return r, err
	}
	return requestHandler(&respBytes)
}

func requestHandler(res *[]byte) (map[string]interface{}, error) {
	var r map[string]interface{}
	err := json.Unmarshal(*res, &r)
	if err != nil {
		fmt.Println(err.Error())
		return r, err
	}
	errCodeTemp, ok := r["errcode"].(float64)
	errCode := int(errCodeTemp)
	if !ok {
		return r, errors.New("响应内容无errorcode")
	}
	//str := (*string)(unsafe.Pointer(&respBytes))
	if errCode == 0 {
		log.Println(string(*res))
		return r, nil
	} else {
		errMsg := r["errmsg"].(string)
		return r, errors.New("接口响应错误码:" + strconv.Itoa(errCode) + "\n错误信息:" + errMsg)
	}
}
