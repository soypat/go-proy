package main

import "fmt"

func main() {
	numba:=3
	var files [numba]string

	file1:= "C:/derpa/herpa"
	file2 := "C:/derpa/sherpe.exe"
	file3 := "C:/derpa/"

	files[0] = file1
	files[1] = file2
	files[2] = file3

	fmt.Println(file1[5:])
	for _,file := range files {
		fmt.Println(file)
	}

}
