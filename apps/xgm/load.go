package main

import (
	_ "embed"
	"fmt"
	"io/fs"
	"peach/data"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"
	"time"
)

// PrintVersion 打印当前数据版本
func PrintVersion(db *sqlite.DB) {
	var ver string
	if err := db.QueryRow("select ver from loadfile where name='xmjh'").Scan(&ver); err == nil {
		fmt.Println("当前数据版本：", ver)
	}
}

// Update 更新计划表
func Update(db *sqlite.DB) (err error) {
	path := utils.NewPath(config.Home).Find("*新柜面存量交易迁移*.xlsx")
	if path == nil {
		return fmt.Errorf("未找到文件：新柜面存量交易迁移*.xlsx")
	}
	fmt.Println("处理文件：", path.Name())
	ver := utils.Extract(`\d{8}`, path.String())
	if f, err := excel.Open(path); err == nil {
		defer f.Close()
		book := f.ExcelBook
		fileinfo := path.FileInfo()
		LoadKfjh(db, fileinfo, &book, ver)
		LoadXjdzb(db, fileinfo, &book, ver)
		LoadXmjh2(db, fileinfo, &book, ver)
	}
	Update_ytc(db)
	update_kfjh(db)
	err = update_kfzt(db)
	Export(db, path)
	return
}

//go:embed query/update_kfjihua.sql
var update_kfjh2 string

func update_kfjh(db *sqlite.DB) {
	fmt.Println("根据验收明细表更新开发状态")
	if r, err := db.Exec(update_kfjh2); err == nil {
		rows, _ := r.RowsAffected()
		utils.Printf("Affected rows: %,d\n", rows)
	} else {
		fmt.Println(err)
	}
}

// LoadXmjh2 项目计划
func LoadXmjh2(db *sqlite.DB, fileinfo fs.FileInfo, book *excel.ExcelBook, ver string) {
	fmt.Println("导入项目计划表")
	if r, err := book.NewReader("全量表", "A:Q", 1, data.HashFilter(-1, -10, -9, -8, -7, -6, -5, -4, -3, -2)); err == nil {
		loader := db.NewLoader(fileinfo, "xmjh", r)
		loader.Ver = ver
		loader.Method = "insert or replace"
		loader.Clear = false
		//loader.Check = false
		//loader.Test(db)
		loader.Load()
	} else {
		fmt.Println(err)
	}
}

// Load 导入数据文件
func Load(db *sqlite.DB) (err error) {
	Home := utils.NewPath(config.Home)
	if path := Home.Find("*新柜面存量交易迁移*.xlsx"); path != nil {
		if err = Load_xmjh(db, path); err != nil {
			fmt.Println(err)
		}
	}
	if path := Home.Find("*数智综合运营系统问题跟踪表*.xlsx"); path != nil {
		if err := LoadWtgzb(db, path); err != nil {
			fmt.Println(err)
		}
	}
	return
}

// Restore 从备份文件中恢复数据
func Restore(db *sqlite.DB) (err error) {
	defer utils.TimeIt(time.Now())
	var path *utils.Path
	if path = utils.NewPath(config.Home).Find("新柜面简报*.zip"); path == nil {
		return fmt.Errorf("未找到 新柜面简报*.zip 文件")
	}
	fmt.Println("处理文件：", path.Name())
	for name, file := range path.IterZip() {
		if strings.Contains(name, "新柜面存量交易迁移计划") {
			err = Load_xmjh(db, file)
		} else if strings.Contains(name, "数智综合运营系统问题跟踪表") {
			err = LoadWtgzb(db, file)
		} else if strings.Contains(name, "版本条目明细") {
			load_bbmx(db, file)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

// Export 更新项目计划表-导出文件
func Export(db *sqlite.DB, path *utils.Path) {
	fmt.Println("更新文件：", path)
	book := excel.NewWriter()
	utils.CheckFatal(ExportAll(db, book, "jh_gbmtj,jh_gzxtj,jh_ywtj,jh_kfjhtj,jh_kfjhb,jh_xmjhb,jh_tcjyb"))
	book.SetColVisible("全量表", "Q", false)
	book.SaveAs(path.String())
	fmt.Println("更新文件完成！")
}
