package main

import (
	"fmt"
	"io"
	"io/fs"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"
)

// conv_gzb 转换故障表的数据
func conv_gzb(src []string) (dest []string, err error) {
	for _, k := range []int{4, 10} {
		src[k] = excel.FormatDate(src[k], "2006-01-02")
	}
	dest = src
	return
}

// LoadWtgzb 导入问题跟踪表数据
func LoadWtgzb(db *sqlite.DB, file utils.File) (err error) {
	var (
		f    io.ReadCloser
		book *excel.ExcelBook
	)
	if f, err = file.Open(); err == nil {
		defer f.Close()
		if book, err = excel.NewExcelBook(f, file.FileInfo().Name()); err == nil {
			load_wtgzb(db, book, file.FileInfo())
		}
	}
	return
}

// load_wtgzb 导入问题跟踪表
func load_wtgzb(db *sqlite.DB, book *excel.ExcelBook, fileinfo fs.FileInfo) (err error) {
	name := fileinfo.Name()
	ver := utils.Extract(`\d{8}`, name)
	fmt.Println("处理文件：", name, "Version:", ver)
	if r, err := book.NewReader(0, "A:M", 1, conv_gzb); err == nil {
		loader := db.NewLoader(fileinfo, "wtgzb", r)
		loader.Ver = strings.Join([]string{ver[:4], ver[4:6], ver[6:]}, "-")
		return loader.Load()
	} else {
		return err
	}
}
