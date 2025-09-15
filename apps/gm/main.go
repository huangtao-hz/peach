package main

import (
	"flag"
	"peach/sqlite"
	"peach/utils"
)

func main() {
	defer utils.Recover()
	db, err := sqlite.Open("xgm")
	utils.CheckFatal(err)
	defer db.Close()
	init_db := flag.Bool("init", false, "初始化数据库")
	load := flag.Bool("load", false, "导入数据")
	query_sql := flag.String("query", "", "执行查询")
	jhbb := flag.String("jhbb", "", "查询计划版本")

	flag.Parse()
	if *init_db {
		CreateDatabse(db)
	}
	if *load {
		//Load(db)
		load_gzb(db)
	}
	if *query_sql != "" {
		db.Println(*query_sql)
	}
	if *jhbb != "" {
		show_jhbb(db, *jhbb)
	}
	if len(flag.Args()) > 0 {
		PrintVersion(db)
	}
	for _, jym := range flag.Args() {
		if utils.FullMatch("\\d{5}", jym) {
			show_new_jy(db, jym)
		} else if utils.FullMatch("\\d{4}", jym) {
			show_old_jy(db, jym)
		}
	}
}
