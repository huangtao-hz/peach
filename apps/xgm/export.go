package main

import (
	_ "embed"
	"fmt"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
)

// export_xjdzb 导出投产交易一览表
func export_xjdzb(db *sqlite.DB, book *excel.Writer) (err error) {
	fmt.Print("导出投产交易一览表，")
	if err := ExportReport(db, book, "tcjyb.toml"); err == nil {
		fmt.Println("完成！")
	}
	return
}

// export_kfjh 导出开发计划表
func export_kfjh(db *sqlite.DB, book *excel.Writer) (err error) {
	ExportReport(db, book, "kfjhb.toml")
	return
}

// export_xmjh 导出项目计划表
func export_xmjh(db *sqlite.DB, book *excel.Writer) (err error) {
	header := "交易码,交易名称,交易组,交易组名,一级菜单,二级菜单,近一年交易量,类型,部门,中心,联系人,方案,计划需求完成时间,当前进度,备注,新交易"
	querys := map[string]string{
		//"计划表": "select * from xmjh where sfwc is null or not sfwc like '5%' order by jym",
		//"完成表": "select * from xmjh where sfwc like '5%' order by jym",
		"全量表": "select *,get_md5(lx,ywbm,zx,lxr,fa,pc,sfwc,bz,xjym) from xmjh order by jym",
	}
	for name, query := range querys {
		sheet := book.GetSheet(name)
		sheet.SetWidth(map[string]float64{
			"A,C":       6.83,
			"B":         42,
			"D":         21,
			"E":         15,
			"F":         31,
			"G,I,J,P,N": 14,
			"H,K,M":     9,
			"O":         24,
			"L":         11,
			"Q":         35,
		})
		sheet.SetColStyle(map[string]string{
			"A:F,H:P": "Normal-NoWrap",
			"G":       "Number",
		})
		sheet.SetColVisible("Q", false)
		if rows, err := db.Query(query); err == nil {
			inch := make(chan []any, BufferSize)
			go rows.FetchAll(inch)
			sheet.AddTable("A1", header, inch)
		} else {
			return err
		}
	}
	return
}

var (
	//go:embed query/tongji.sql
	tongji_sql string
	//go:embed query/tongji_gzx.sql
	tongji_gzx_sql string
	//go:embed query/tongji_gbm.sql
	tongji_gbm_sql string

	//go:embed tables/tjbtable.toml
	table_tjb string
	//go:embed tables/gbmtable.toml
	table_gbm string
	//go:embed tables/gzxtable.toml
	table_gzx string
	//go:embed tables/kfjhtj.toml
	table_kfjhtj string
)

// export_tjb 导出统计表
func export_tjb(db *sqlite.DB, book *excel.Writer) (err error) {
	sheet := book.GetSheet("统计表")
	sheet.SetWidth(map[string]float64{
		"A":   12,
		"B":   20,
		"C:H": 10,
	})

	rows, err := db.Query(tongji_gbm_sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.AddTableToml("B1", table_gbm, ch)

	//header = "中心,未提交需求,已完成需求,开发中,已投产,总数,投产完成率"
	rows, err = db.Query(tongji_gzx_sql)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch = make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	//sheet.AddTable("B8", header, ch)
	sheet.AddTableToml("B8", table_gzx, ch)

	//header = "联系人,中心,未提交需求,已完成需求,开发中,已投产,总数,投产完成率"
	rows, err = db.Query(tongji_sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch = make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	//sheet.AddTable("A19", header, ch)
	err = sheet.AddTableToml("A19", table_tjb, ch)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(table_tjb)
	return
}

// export_kfjhtj 导出开发计划统计表
func export_kfjhtj(db *sqlite.DB, book *excel.Writer) (err error) {
	sheet := book.GetSheet("开发计划统计")
	sheet.SetWidth(map[string]float64{
		"A:C": 12,
	})
	query := `
select a.jhbb,count(a.jym),count(a.jym)*1.0/(select count(jym)from xmjh where fa not in ("1-下架交易","5-移出柜面系统"))
from kfjh a
left join xmjh b
on a.jym=b.jym
where b.sfwc not like "5%"
group by jhbb
order by jhbb
	`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.AddTableToml("A1", table_kfjhtj, ch)
	return
}

// Export 更新项目计划表-导出文件
func Export(db *sqlite.DB, path *utils.Path) (err error) {
	fmt.Println("更新文件：", path)
	book := excel.NewWriter()
	export_tjb(db, book)
	export_kfjhtj(db, book)
	export_kfjh(db, book)
	export_xmjh(db, book)
	export_xjdzb(db, book)
	book.SaveAs(path.String())
	fmt.Println("更新文件完成！")
	return
}
