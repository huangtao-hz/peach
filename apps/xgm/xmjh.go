package main

import (
	"fmt"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"
)

// load_xmjh 导入项目计划
func load_xmjh(db *sqlite.DB, file utils.File) (err error) {
	var f *excel.ExcelFile
	fileinfo := file.FileInfo()
	date := utils.Extract(`\d{8}`, fileinfo.Name())
	date = strings.Join([]string{date[:4], date[4:6], date[6:]}, "-")
	if f, err = excel.Open(file); err == nil {
		fmt.Println("导入新旧交易对照表")
		db.LoadExcel(loaderFS, "loader/jh_xjdzb.toml", &f.ExcelBook, fileinfo, date)
		fmt.Println("导入开发计划表")
		db.LoadExcel(loaderFS, "loader/jh_kfjh.toml", &f.ExcelBook, fileinfo, date)
		fmt.Println("导入项目计划表")
		db.LoadExcel(loaderFS, "loader/jh_xmjh.toml", &f.ExcelBook, fileinfo, date)
		fmt.Println("导入版本安排")
		db.LoadExcel(loaderFS, "loader/jh_bbap.toml", &f.ExcelBook, fileinfo, date)
	}
	return
}

// export_xmjh 更新项目计划表-导出文件
func export_xmjh(db *sqlite.DB, path *utils.Path) {
	fmt.Print("更新文件：", path.Base())
	utils.CheckFatal(ExportXlsx(db, path.String(), "jh_gbmtj,jh_gzxtj,jh_ywtj,jh_kfjhtj,jh_kfjhhztj,jh_gzxkfjh,jh_kfjhb,jh_xmjhb,jh_tcjyb,jh_bbap"))
	fmt.Println(" 完成！")
}

// PrintVersion 打印当前数据版本
func PrintVersion(db *sqlite.DB) {
	var ver string
	if err := db.QueryRow("select ver from loadfile where name='xmjh'").Scan(&ver); err == nil {
		fmt.Println("当前数据版本：", ver)
	}
}
