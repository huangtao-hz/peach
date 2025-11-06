package main

import (
	"fmt"
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
	if path != nil {
		fmt.Println("处理文件：", path.Name())
		if f, err := excel.Open(path); err == nil {
			defer f.Close()
			fileinfo := path.FileInfo()
			fmt.Println("导入新旧交易对照表")
			db.LoadExcel(loaderFS, "loader/jh_xjdzb.toml", &f.ExcelBook, fileinfo)
			fmt.Println("导入开发计划表")
			db.LoadExcel(loaderFS, "loader/jh_kfjh.toml", &f.ExcelBook, fileinfo)
			fmt.Println("导入项目计划表")
			db.LoadExcel(loaderFS, "loader/jh_xmjh2.toml", &f.ExcelBook, fileinfo, data.HashFilter(-1, -10, -9, -8, -7, -6, -5, -4, -3, -2))
			fmt.Println("导入版本安排")
			db.LoadExcel(loaderFS, "loader/jh_bbap.toml", &f.ExcelBook, fileinfo)
		}
	} else {
		path = utils.NewPath(config.Home).Join(fmt.Sprintf("附件1：新柜面存量交易迁移计划%s.xlsx", utils.Today().Format("%Y%M%D")))
	}
	//load_kfjh(db) // 导入科技管理部编制的开发计划表
	fmt.Print("更新验收完成时间:")
	r, _ := db.Exec(`update bbap set wcys=date(tcrq,"weekday 5","-7 days") where wcys=""`)
	if count, err := r.RowsAffected(); err == nil {
		fmt.Println(count, "条数据被更新")
	}
	Update_ytc(db)
	update_kfjh(db)
	fmt.Print("根据计划版本更新开发计划时间：")
	db.ExecuteFs(queryFS, "query/update_kfjhsj.sql")
	fmt.Print("根据验收条目更新完成状态：")
	db.ExecuteFs(queryFS, "query/update_xmjh.sql")
	fmt.Print("根据新旧交易对照表更新对应新交易：")
	db.ExecuteFs(queryFS, "query/update_xmjh_xjy.sql")
	Export(db, path)
	return
}

// Export 更新项目计划表-导出文件
func Export(db *sqlite.DB, path *utils.Path) {
	fmt.Println("更新文件：", path)
	utils.CheckFatal(ExportXlsx(db, path.String(), "jh_gbmtj,jh_gzxtj,jh_ywtj,jh_kfjhtj,jh_gzxkfjh,jh_kfjhb,jh_xmjhb,jh_tcjyb,jh_bbap"))
	fmt.Println("更新文件完成！")
}

func update_kfjh(db *sqlite.DB) {
	fmt.Print("根据验收明细表更新开发状态:")
	db.ExecuteFs(queryFS, "query/update_kfjihua.sql")
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
