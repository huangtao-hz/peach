package main

import (
	"fmt"
	"peach/sqlite"
	"peach/utils"
	"strings"
)

func query(db *sqlite.DB, args ...string) {
	if len(args) == 1 && utils.FullMatch(`\d{4,5}`, args[0]) {
		db.Println("select jym,jymc from bbsm where jym=? order by rq desc limit 1", args[0])
		fmt.Println()
		db.Println("select printf('%s    %s\n%s\n',rq,lxr,nr) from bbsm where jym=? order by rq", args[0])
	} else {
		query := make([]string, len(args))
		for i, arg := range args {
			query[i] = fmt.Sprintf("nr like '%%%s%%'", arg)
		}
		aquery := strings.Join(query, " and ")
		db.Println(fmt.Sprintf("select printf('%%s-%%s\n%%s\n%%s by:%%s\n',jym,jymc,nr,rq,lxr) from bbsm  where %s", aquery))

	}
}

func show_date(db *sqlite.DB) {
	header := "投产日期    优化数量"
	query := "select rq,sum(sl)from bbsm_view group by rq order by rq"
	db.Printf(query, "%s      %5,d\n", header, true)
}

func show_year(db *sqlite.DB) {
	header := "年份    投产次数   优化数量"
	query := "select strftime('%Y',rq)as nf,count(distinct rq),sum(sl)from bbsm_view group by nf order by nf"
	db.Printf(query, "%s       %4,d       %5,d\n", header, true)
}
