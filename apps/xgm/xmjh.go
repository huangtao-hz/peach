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

// Export 更新项目计划表-导出文件
func Export(db *sqlite.DB, path *utils.Path) {
	fmt.Print("更新文件：", path.Base())
	utils.CheckFatal(ExportXlsx(db, path.String(), "jh_gbmtj,jh_gzxtj,jh_ywtj,jh_kfjhtj,jh_gzxkfjh,jh_kfjhb,jh_xmjhb,jh_tcjyb,jh_bbap"))
	fmt.Println(" 完成！")
}

func Update_ytc(db *sqlite.DB) {
	var count int
	err := db.QueryRow(`select count(a.jym) from xmjh a join xjdz b on a.jym=b.yjym and b.tcrq<date('now') and b.tcrq<>"" and a.sfwc not like '5%' `).Scan(&count)
	utils.CheckFatal(err)
	if count > 0 {
		utils.Printf("以下交易已有新旧交易对照表，共有%,d条记录。\n", count)
		db.Println(`select a.jym,a.jymc,a.sfwc,b.jym,b.jymc,b.tcrq from xmjh a join xjdz b on a.jym=b.yjym
			where b.tcrq<date('now') and a.sfwc not like '5%' and b.tcrq<>''`)
		sql := `update xmjh
                set sfwc='5-已投产'
                from xjdz
                where sfwc<>'5-已投产' and xmjh.jym=xjdz.yjym and tcrq<date('now') and tcrq<>'' `
		fmt.Print("请确认是否修改？Y or N:")
		var s string
		if n, err := fmt.Scan(&s); err == nil && n > 0 && strings.ToUpper(s) == "Y" {
			if r, err := db.Exec(sql); err == nil {
				if count, err := r.RowsAffected(); err == nil {
					utils.Printf("%,d 条数据被修改 \n", count)
				}
			} else {
				fmt.Println(err)
			}
		}
	}
}

// PrintVersion 打印当前数据版本
func PrintVersion(db *sqlite.DB) {
	var ver string
	if err := db.QueryRow("select ver from loadfile where name='xmjh'").Scan(&ver); err == nil {
		fmt.Println("当前数据版本：", ver)
	}
}
