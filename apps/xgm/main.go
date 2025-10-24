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
var Version = "1.0.2"

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
	init_db := flag.Bool("init", false, "初始化数据库")
	load := flag.Bool("load", false, "导入数据")
	query_sql := flag.String("query", "", "执行查询")
	jhbb := flag.String("jhbb", "", "查询计划版本")
	restore := flag.Bool("restore", false, "导入数据")
	touchan := flag.Bool("touchan", false, "导入数据")
	update := flag.Bool("update", false, "更新计划表")
	jihua := flag.Bool("jihua", false, "投产交易清单")
	version := flag.Bool("version", false, "查阅程序版本")

	flag.Parse()
	if *init_db {
		fmt.Println("初始化数据库表")
		CreateDatabse(db)
		fmt.Println("初始化数据库成功！")
	}
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

		err = Update(db)
		//err = Export(db, nil)
		update_bbmx(db)
	}
	if *jihua {
		kaifajihua(db)
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
	if err != nil {
		fmt.Println("Error:", err)
	}
}
