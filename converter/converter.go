package converter

import (
	"reflect"
	"strings"
)

var byteMap = map[string]byte{
	"_":  0x5F,
	"A":  0x41,
	"Z":  0x5A,
	"a":  0x61,
	"z":  0x7A,
	"32": 0x20,
}

func HumpToSnake(str string) string {
	bytes := []byte(str)
	newBytes := []byte("")
	for i, char := range bytes {
		//上一个是否为大写
		lastIsUpper := false
		if i > 0 {
			lastChar := bytes[i-1]
			if lastChar >= byteMap["A"] && lastChar <= byteMap["Z"] {
				lastIsUpper = true
			}
		}
		//小写数字符号的情况
		temp := []byte{char}
		//大写的情况
		if char >= byteMap["A"] && char <= byteMap["Z"] {
			temp = []byte{char + byteMap["32"]}
			//第一个、上一个是大写 都不处理
			if i != 0 && lastIsUpper == false && bytes[i] != byteMap["_"] {
				temp = []byte{byteMap["_"], char + byteMap["32"]}
			}
		}
		newBytes = append(newBytes, temp...)
	}
	return string(newBytes)
}

func SnakeToHump(str string, firstUpper bool) string {
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, "___", "_")
	str = strings.ReplaceAll(str, "__", "_")
	bytes := []byte(str)
	newBytes := []byte("")
	length := len(bytes)
	for i := 0; i < length; i++ {
		char := bytes[i]
		if i == 0 && char == byteMap["_"] {
			continue
		}
		temp := []byte{char}
		if char == byteMap["_"] {
			if i+1 >= length {
				continue
			}
			nextChar := bytes[i+1]
			if nextChar >= byteMap["a"] && nextChar <= byteMap["z"] {
				temp = []byte{nextChar - byteMap["32"]}
			} else {
				temp = []byte{nextChar}
			}
			i += 1
		}
		newBytes = append(newBytes, temp...)
	}
	if firstUpper {
		newBytes = append([]byte{newBytes[0] - byteMap["32"]}, newBytes[1:]...)
	}
	return string(newBytes)
}

func StructFieldsNameArray(instance interface{}, useSnakeName bool) (result []string) {
	result = []string{}
	//获取一组反射值
	ref := reflect.ValueOf(instance)
	//遍历这狙数据
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Type().Field(i).Name
		if useSnakeName {
			field = HumpToSnake(field)
		}
		result = append(result,field)
	}
	return result
}
