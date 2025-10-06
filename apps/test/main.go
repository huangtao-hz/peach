package main

import (
	"fmt"
	"peach/archive"
	"peach/excel"
	"peach/utils"
)

func proc(file archive.File) {
	if !file.FileInfo().IsDir() {
		fmt.Println("File:", file.FileInfo().Name())
		if book, err := excel.OpenFile(file); err == nil {
			fmt.Println(book.GetSheetList())
		} else {
			fmt.Println("Error:", err)
		}
	}
}

func main() {
	defer utils.Recover()
	file := utils.NewPath("~/Downloads").Find("20250905.tgz")
	if file != nil {
		fmt.Println(file.String())
		//fmt.Println(file.FileInfo().Name())
		archive.ExtractTar(file.String(), proc)
	}
}
