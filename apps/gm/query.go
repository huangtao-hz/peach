package main

import (
	"fmt"
	"peach/sqlite"
)

const (
	xjyHeader = "交易码,交易名称,投产日期,状态,备注"
	XjyQuery  = `select b.jym,b.jymc,b.ywbm,b.zx,b.lxr,a.tcrq
	from xjdz a
	left join xmjh b
	on a.yjym=b.jym
	where a.jym=?`

	XyjOldHeader = `交易码       交易名称                 业务部门       中心               联系人   投产日期`
)

func show_new_jy(db *sqlite.DB, jym string) {
	db.PrintRow("select jym,jymc,min(tcrq),zs,bz from xjdz where jym=? group by jym", xjyHeader, jym)
	fmt.Println("                      --  对应老交易清单  --")
	db.Printf(XjyQuery, "%4s  %-30s %-12s  %-12s  %12s    %10s\n\n", XyjOldHeader, true, jym)
}

func show_old_jy(db *sqlite.DB, jym string) {
	var fa, zt, jhbb string
	err := db.QueryRow("select fa,sfwc,ifnull(jhbb,'') from xmjh a left join kfjh b on a.jym=b.jym where a.jym=?", jym).Scan(&fa, &zt, &jhbb)
	if err != nil {
		fmt.Printf("交易 %s 不存在\n", jym)
		return
	}
	if fa[0] == '1' || fa[0] == '5' {
		header := "交易码,交易名称,交易笔数,类型,业务部门,中心,业务联系人,改造方案"
		sql := "select a.jym,a.jymc,a.bs,a.lx,a.ywbm,a.zx,a.lxr,a.fa from xmjh a left join kfjh b on a.jym=b.jym where a.jym=?"
		db.PrintRow(sql, header, jym)
	} else if zt[0] == '5' {
		header := "交易码,交易名称,交易笔数,类型,业务部门,中心,业务联系人,改造方案,状态,对应新交易,投产日期"
		sql := `select a.jym,a.jymc,a.bs,a.lx,a.ywbm,a.zx,a.lxr,a.fa,a.sfwc,printf('%s-%s',c.jym,c.jymc),c.tcrq
    from xmjh a
    left join xjdz c on a.jym=c.yjym
    where a.jym=?`
		db.PrintRow(sql, header, jym)
	} else if jhbb == "" {
		header := "交易码,交易名称,交易笔数,类型,业务部门,中心,业务联系人,改造方案,状态"
		sql := `select a.jym,a.jymc,a.bs,a.lx,a.ywbm,a.zx,a.lxr,a.fa,printf("%s（未制定计划）",a.sfwc) from xmjh a where a.jym=?`
		db.PrintRow(sql, header, jym)
	} else {
		header := "交易码,交易名称,交易笔数,类型,业务部门,中心,业务联系人,改造方案,计划版本,技术经理,开发组长"
		sql := "select a.jym,a.jymc,a.bs,a.lx,a.ywbm,a.zx,a.lxr,a.fa,b.jhbb,b.kjfzr,b.kfzz from xmjh a left join kfjh b on a.jym=b.jym where a.jym=?"
		db.PrintRow(sql, header, jym)
	}
	fmt.Println("")
}

// show_jhbb 查询指定版本的交易明细
func show_jhbb(db *sqlite.DB, jhbb string) {
	const (
		query = `select a.jym,b.jymc,b.ywbm,b.zx,b.lxr,a.kjfzr,a.kfzz
from kfjh a left join xmjh b on a.jym=b.jym
where a.jhbb=?
order by a.jym`
		header = "交易码  交易名称           部门       中心       联系人    技术经理        开发组长"
		format = "%4s %-35s %-12s %-16s %-10s %-10s %-10s\n"
	)
	db.Printf(query, format, header, true, jhbb)
}

// show_tongji 打印各版本的统计信息
func show_tongji(db *sqlite.DB) {
	header := "投产日期   交易数量 迁移交易数量 新交易数量    占比（%）"
	sql := `select tcrq,count(distinct jym),sum(iif(yjym<>"",1,0)),sum(iif(yjym="",1,0)),
    sum(iif(yjym<>"",1,0))*100.0/(select count(jym)from xmjh where fa not in ("1-下架交易","5-移出柜面系统"))
    from xjdz
    where tcrq<=date('now')
    group by tcrq
    union
    -- 显示合计数据
    select '合计',count(distinct jym),sum(iif(yjym<>"",1,0)),sum(iif(yjym="",1,0)),
    sum(iif(yjym<>"",1,0))*100.0/(select count(jym)from xmjh where fa not in ("1-下架交易","5-移出柜面系统"))
    from xjdz
    where tcrq<=date('now')
`
	format := "%10s  %8,d  %8,d  %8,d        %5.2f %%\n"
	db.Printf(sql, format, header, true)
}
