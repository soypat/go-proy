package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var files []string
	root,err:=filepath.Abs("./")
	if err!=nil{
		return
	}


	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		fmt.Println(files)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
