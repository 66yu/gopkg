package converter

import (
	"testing"
	"time"
)

func TestHumpToSnack(t *testing.T) {
	var originSet = []string{"MySql", "WordConverter","myWord","HelloWorld","ID"}
	var resultSet = []string{"my_sql", "word_converter","my_word","hello_world","id"}
	for i, str := range originSet {
		if resultSet[i] != HumpToSnake(str) {
			t.Error("HumpToSnake error ", str, HumpToSnake(str))
		}
	}
}

func TestSnakeToHump(t *testing.T) {
	var originSet = []string{"my_sql", "word_converter","my_word","______heLlo_wORLd_","_hello_world_","sqlite_3_driver_"}
	var resultSet = []string{"MySql", "WordConverter","myWord","HelloWorld","Sqlite3Driver"}
	for _, str := range originSet {
		def := SnakeToHump(str, false)
		firstUpper := SnakeToHump(str, true)
		pass := false
		for _, result := range resultSet {
			if result == def || result == firstUpper {
				pass = true
			}
		}
		if pass==false {
			t.Error("SnakeToHump error ", str, def, firstUpper)
		}
	}
}

func TestStructFieldsNameArray(t *testing.T) {
	type test struct {
		ID int64
		Name string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt time.Time
	}
	result := []string{"id"}
	arr := StructFieldsNameArray(test{},true)
	for _, s := range result {
		exist := false
		for _, s2 := range arr {
			if s==s2 {
				exist = true
			}
		}
		if exist!=true {
			t.Error("SnakeToHump error: ", s)
		}
	}
}
