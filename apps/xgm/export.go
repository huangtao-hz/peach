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
	ExportReport(db, book, "xmjhb.toml")
	book.SetColVisible("全量表", "Q", false)
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
	return ExportReport(db, book, "kfjhtj.toml")
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
