package main

import (
	_ "embed"
	"flag"
	"peach/sqlite"
	"peach/utils"
)

//go:embed db.sql
var create_sql string

func main() {
	defer utils.Recover()
	db, err := sqlite.Open("bbsm")
	utils.CheckFatal(err)
	defer db.Close()
	init := flag.Bool("init", false, "初始化数据库")
	load := flag.Bool("load", false, "导入数据")
	date := flag.Bool("date", false, "按日期统计")
	year := flag.Bool("year", false, "按年统计")
	flag.Parse()
	if *init {
		db.ExecScript(create_sql)
		sqlite.InitLoadFile(db)
	}
	if *load {
		load_all(db)
	}
	if *date {
		show_date(db)
	}
	if *year {
		show_year(db)
	}
	if len(flag.Args()) > 0 {
		query(db, flag.Args()...)
	}

}
