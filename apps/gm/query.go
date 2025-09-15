package main

import (
	"fmt"
	"peach/sqlite"
)

const (
	xjyHeader = "交易码,交易名称,投产日期,状态,备注"
	XjyQuery  = `select b.jym,b.jymc,b.bm,b.zx,b.lxr
	from jydzb a
	left join xmjh b
	on a.jym=b.jym
	where a.xjym=?`

	XyjOldHeader = `交易码       交易名称                 业务部门      中心          联系人`
)

func show_new_jy(db *sqlite.DB, jym string) {
	db.PrintRow("select distinct xjym,xjymc,tcrq,zt,bz from jydzb where xjym=?", xjyHeader, jym)
	fmt.Println("                      --  对应老交易清单  --")
	db.Printf(XjyQuery, "%4s  %-30s %-12s  %-12s  %12s\n", XyjOldHeader, true, jym)
}

func show_old_jy(db *sqlite.DB, jym string) {
	const (
		LjyHeader = "交易码,交易名称,交易组,菜单,业务部门,联系人,方案,业务量,投产版本,技术经理,开发组长,对应新交易,投产日期"
		LjyQuery  = `select a.jym,a.jymc,printf('%s-%s',a.jyz,a.jyzm),printf('%s -> %s',a.yjcd,a.ejcd),
printf('%s-%s',a.bm,a.zx),a.lxr,a.fa,a.ywl,
ifnull(b.jhbb,''),ifnull(b.kffzr,''),ifnull(b.kfzz,''),
printf("%s-%s",c.xjym,c.xjymc),ifnull(c.tcrq,'')
from xmjh a
left join kfjh b on a.jym=b.jym
left join jydzb c on a.jym=c.jym
where a.jym=?`
	)
	db.PrintRow(LjyQuery, LjyHeader, jym)
}

func show_jhbb(db *sqlite.DB, jhbb string) {
	const (
		query = `select a.jym,b.jymc,b.bm,b.zx,b.lxr,a.kffzr,a.kfzz
from kfjh a left join xmjh b on a.jym=b.jym
where a.jhbb=?
order by a.jym`
		header = "交易码  交易名称           部门       中心       联系人    技术经理        开发组长"
		format = "%4s %-35s %-12s %-16s %-10s %-10s %-10s\n"
	)
	db.Printf(query, format, header, true, jhbb)
}
