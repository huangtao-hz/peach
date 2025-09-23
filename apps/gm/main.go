package main

import (
	"flag"
	"fmt"
	"peach/sqlite"
	"peach/utils"
	"strings"
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
	fix := flag.Bool("fix", false, "修正数据")

	flag.Parse()
	if *init_db {
		CreateDatabse(db)
	}
	if *fix {
		var count int
		err := db.QueryRow(`select count(a.jym) from xmjh a join xjdz b on a.jym=b.yjym and b.tcrq<date('now') and a.sfwc not like '5%' `).Scan(&count)
		utils.CheckFatal(err)
		if count > 0 {
			utils.Printf("以下交易已有新旧交易对照表，共有%,d条记录。\n", count)
			db.Println(`select a.jym,a.jymc,a.sfwc,b.jym,b.jymc,b.tcrq from xmjh a join xjdz b on a.jym=b.yjym
				where b.tcrq<date('now') and a.sfwc not like '5%' `)
			sql := `update xmjh
                    set sfwc='5-已投产'
                    from xjdz
                    where sfwc<>'5-已投产' and xmjh.jym=xjdz.yjym and xjdz.tcrq <date('now')`
			fmt.Print("请确认是否修改？Y or N:")
			var s string
			if n, err := fmt.Scan(&s); err == nil && n > 0 && strings.ToUpper(s) == "Y" {
				if r, err := db.Exec(sql); err == nil {
					if count, err := r.RowsAffected(); err == nil {
						utils.Printf("%,d 条数据被修改。\n", count)
					}
				} else {
					fmt.Println(err)
				}
			}
		}
	}
	if *load {
		Load(db)
		//load_gzb(db)
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
		if utils.FullMatch("\\d{5}", jym) {
			show_new_jy(db, jym)
		} else if utils.FullMatch("\\d{4}", jym) {
			show_old_jy(db, jym)
		}
	}
}
