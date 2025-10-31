package main

import (
	"fmt"
	"io"
	"io/fs"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
)

// Load_xmjh 导入项目计划
func Load_xmjh(db *sqlite.DB, file utils.File) (err error) {
	var (
		f    io.ReadCloser
		book *excel.ExcelBook
	)
	if f, err = file.Open(); err == nil {
		defer f.Close()
		fileinfo := file.FileInfo()
		name := fileinfo.Name()
		ver := utils.Extract(`\d{8}`, name)
		fmt.Println("处理文件：", name, "Version:", ver)
		if book, err = excel.NewExcelBook(f, name); err == nil {
			utils.CheckErr(db.LoadExcel(loaderFS, "loader/jh_xmjh.toml", book, fileinfo))
			utils.CheckErr(db.LoadExcel(loaderFS, "loader/jh_xjdzb.toml", book, fileinfo))
			//LoadXjdzb(db, fileinfo, book, ver)
			//LoadXmjh(db, fileinfo, book, ver)
		}
		Update_ytc(db)
	}
	return
}

// LoadKfjh 开发计划
func LoadKfjh(db *sqlite.DB, fileinfo fs.FileInfo, book *excel.ExcelBook, ver string) {
	fmt.Println("导入开发计划表")
	if r, err := book.NewReader("开发计划", "A,H:S", 1); err == nil {
		loader := db.NewLoader(fileinfo, "kfjh", r)
		loader.Ver = ver
		//loader.Check = false
		loader.Load()
	} else {
		fmt.Println(err)
	}
}
