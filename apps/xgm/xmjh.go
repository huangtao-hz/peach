package main

import (
	"fmt"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
)

// load_xmjh 导入项目计划
func load_xmjh(db *sqlite.DB, file utils.File) (err error) {
	var f *excel.ExcelFile
	fileinfo := file.FileInfo()
	if f, err = excel.Open(file); err == nil {
		fmt.Println("导入新旧交易对照表")
		db.LoadExcel(loaderFS, "loader/jh_xjdzb.toml", &f.ExcelBook, fileinfo)
		fmt.Println("导入开发计划表")
		db.LoadExcel(loaderFS, "loader/jh_kfjh.toml", &f.ExcelBook, fileinfo)
		fmt.Println("导入项目计划表")
		db.LoadExcel(loaderFS, "loader/jh_xmjh.toml", &f.ExcelBook, fileinfo)
		fmt.Println("导入版本安排")
		db.LoadExcel(loaderFS, "loader/jh_bbap.toml", &f.ExcelBook, fileinfo)
	}
	return
}

// Export 更新项目计划表-导出文件
func Export(db *sqlite.DB, path *utils.Path) {
	fmt.Println("更新文件：", path)
	utils.CheckFatal(ExportXlsx(db, path.String(), "jh_gbmtj,jh_gzxtj,jh_ywtj,jh_kfjhtj,jh_gzxkfjh,jh_kfjhb,jh_xmjhb,jh_tcjyb,jh_bbap"))
	fmt.Println("更新文件完成！")
}
