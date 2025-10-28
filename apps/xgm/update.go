package main

import (
	"fmt"
	"peach/sqlite"
	"peach/utils"
	"strings"
)

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
					utils.Printf("%,d 条数据被修改。\n", count)
				}
			} else {
				fmt.Println(err)
			}
		}
	}
}
