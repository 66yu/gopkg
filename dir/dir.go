package dir

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func Mkdir(path string,perm fs.FileMode) error  {
	log.Print(filepath.Abs(path))
	if _,err := os.Stat(path);err!=nil{
		err = os.MkdirAll(path,perm)
		if err != nil {
			return err
		}
	}
	return nil
}