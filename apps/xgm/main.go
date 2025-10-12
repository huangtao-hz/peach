package main

import (
	_ "embed"
	"flag"
	"fmt"
	"peach/sqlite"
	"peach/utils"
)

//go:embed xmjh.sql
var create_sql string

func CreateDatabse(db *sqlite.DB) {
	fmt.Println("初始化数据库表")
	db.ExecScript(create_sql)
	sqlite.InitLoadFile(db)
	fmt.Println("初始化数据库成功！")
}

// main 主程序入口
func main() {
	//defer utils.Recover()
	db, err := sqlite.Open("xgm2025-03")
	utils.CheckFatal(err)
	defer db.Close()
	init_db := flag.Bool("init", false, "初始化数据库")
	load := flag.Bool("load", false, "导入数据")
	query_sql := flag.String("query", "", "执行查询")
	jhbb := flag.String("jhbb", "", "查询计划版本")
	restore := flag.Bool("restore", false, "导入数据")
	tongji := flag.Bool("tongji", false, "导入数据")
	update := flag.Bool("update", false, "更新计划表")

	flag.Parse()
	if *init_db {
		CreateDatabse(db)
	}
	if *load {
		err = load_jyjh(db)
		//LoadWtgzb(db)
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
	if *update {
		err = Update(db)
		//err = Export(db, nil)
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
	/*
		fmt.Println(s)
		table := &excelize.Table{}
		err = toml.Unmarshal([]byte(s), table)
		fmt.Println(err)
		utils.PrintStruct(table)
	*/
}
