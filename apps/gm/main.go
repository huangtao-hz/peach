package main

import (
	"flag"
	"peach/sqlite"
	"peach/utils"
)

// main 主程序入口
func main() {
	defer utils.Recover()
	db, err := sqlite.Open("xgm2025-03")
	utils.CheckFatal(err)
	defer db.Close()
	init_db := flag.Bool("init", false, "初始化数据库")
	load := flag.Bool("load", false, "导入数据")
	query_sql := flag.String("query", "", "执行查询")
	jhbb := flag.String("jhbb", "", "查询计划版本")
	restore := flag.Bool("restore", false, "导入数据")
	tongji := flag.Bool("tongji", false, "导入数据")

	flag.Parse()
	if *init_db {
		CreateDatabse(db)
	}
	if *load {
		//Load(db)
		LoadWtgzb(db)
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
	if *tongji {
		show_tongji(db)
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
}
