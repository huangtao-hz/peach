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

// update_xmjh 更新计划表
func update_xmjh(db *sqlite.DB) (err error) {
	path := utils.NewPath(config.Home).Find("*新柜面存量交易迁移*.xlsx")
	if path != nil {
		load_xmjh(db, path)
	} else {
		path = utils.NewPath(config.Home).Join(fmt.Sprintf("附件1：新柜面存量交易迁移计划%s.xlsx", utils.Today().Format("%Y%M%D")))
	}
	//load_kfjh(db) // 导入科技管理部编制的开发计划表
	fmt.Print("根据投产时间更新验收完成时间:")
	r, _ := db.Exec(`update bbap set wcys=date(tcrq,"weekday 5","-7 days") where wcys=""`)
	if count, err := r.RowsAffected(); err == nil {
		fmt.Println(count, "条数据被更新")
	}
	fmt.Print("根据验收明细表更新开发状态:")
	db.ExecuteFs(queryFS, "query/update_kfjihua.sql")
	fmt.Print("根据计划版本更新开发计划时间：")
	db.ExecuteFs(queryFS, "query/update_kfjhsj.sql")
	fmt.Print("根据验收条目更新完成状态：")
	db.ExecuteFs(queryFS, "query/update_xmjh.sql")
	fmt.Print("根据新旧交易对照表更新对应新交易：")
	db.ExecuteFs(queryFS, "query/update_xmjh_xjy.sql")
	fmt.Print("更新当前版本已投产交易：")
	db.ExecuteFs(queryFS, "query/update_xmjh_ytc.sql")
	export_xmjh(db, path)
	return
}
