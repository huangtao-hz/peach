package main

import (
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
)

// Load_xmjh 导入项目计划
func Load_xmjh(db *sqlite.DB, file utils.File) (err error) {
	var f *excel.ExcelFile
	fileinfo := file.FileInfo()
	if f, err = excel.Open(file); err == nil {
		utils.CheckErr(db.LoadExcel(loaderFS, "loader/jh_xmjh.toml", &f.ExcelBook, fileinfo))
		utils.CheckErr(db.LoadExcel(loaderFS, "loader/jh_xjdzb.toml", &f.ExcelBook, fileinfo))
		//utils.CheckErr(db.LoadExcel(loaderFS, "loader/jh_kfjh.toml", &f.ExcelBook, fileinfo))
	}
	return
}
