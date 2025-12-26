package main

import (
	_ "embed"
	"flag"
	"fmt"
	"peach/sqlite"
	"peach/utils"
)

//go:embed query/db.sql
var create_sql string
var Version = "1.0.7"

//go:embed query/update_ysrq.sql
var update_ysrq string

func CreateDatabse(db *sqlite.DB) {
	db.ExecScript(create_sql)
	sqlite.InitLoadFile(db)
}

// main 主程序入口
func main() {
	defer utils.Recover()
	db, err := sqlite.Open(config.Database)
	utils.CheckFatal(err)
	defer db.Close()
	load := flag.Bool("load", false, "导入数据")
	query_sql := flag.String("query", "", "执行查询")
	jhbb := flag.String("jhbb", "", "查询计划版本")
	restore := flag.Bool("restore", false, "导入数据")
	touchan := flag.Bool("touchan", false, "导入数据")
	update := flag.Bool("update", false, "更新计划表")
	jihua := flag.Bool("jihua", false, "投产交易清单")
	version := flag.Bool("version", false, "查阅程序版本")
	wenti := flag.String("wenti", "", "统计上报问题，取值：本月、上月、上周、本周")

	flag.Parse()
	CreateDatabse(db)
	if *version {
		fmt.Println("版本：", Version)
	}
	if *load {
		//err = load_jyjh(db)
		err = Load(db)
	}
	if *query_sql != "" {
		db.Println(*query_sql)
	}
	if *jhbb != "" {
		show_jhbb(db, *jhbb)
	}
	if *restore {
		Restore(db)
	}
	if *touchan {
		show_touchan(db)
	}
	if *update {
		load_qxzb(db)
		update_bbmx(db)
		err = Update(db)
		Home := utils.NewPath(config.Home)
		if path := Home.Find("*数智综合运营系统问题跟踪表*.xlsx"); path != nil {
			load_wtgzb(db, path)
		}
	}
	if *jihua {
		kaifajihua(db)
	}
	if *wenti != "" {
		report_wenti(db, *wenti)
	}
	if len(flag.Args()) > 0 {
		PrintVersion(db)
	}
	for _, jym := range flag.Args() {
		if utils.FullMatch(`\d{5}`, jym) {
			show_new_jy(db, jym)
		} else if utils.FullMatch(`\d{4}`, jym) {
			show_old_jy(db, jym)
		}
	}
	utils.CheckErr(err)
}
