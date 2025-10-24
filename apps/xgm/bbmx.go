package main

import (
	_ "embed"
	"fmt"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"
)

// conv_jym 转换交易码
func conv_jym(jym *string) {
	if len(*jym) == 3 {
		*jym = fmt.Sprintf("%04s", *jym)
	}
}

// load_bbsm 导入版本明细数据
func load_bbmx(db *sqlite.DB, path *utils.Path) {
	if f, err := excel.NewExcelFile(path.String()); err == nil {
		defer f.Close()
		fmt.Println("导入需求条目表")
		date := utils.Extract(`\d{8}`, path.String())
		date = strings.Join([]string{date[:4], date[4:6], date[6:]}, "-")
		var conv = func(src []string) (dest []string, err error) {
			if src[0] == "" {
				return
			}
			dest = []string{date}
			dest = append(dest, src...)
			return
		}
		ystm_reader, _ := f.NewReader("需求条目列表", "A:K", 1, conv)
		loader := db.NewLoader(path.FileInfo(), "ystmb", ystm_reader)
		loader.Load()

		fmt.Println("导入新旧对照表")
		var conv_jydz = func(src []string) (dest []string, err error) {
			if src[0] == "" {
				return
			}
			conv_jym(&src[1])
			dest = src
			return
		}
		jydz_reader, _ := f.NewReader("新旧交易对照表", "A,C,D", 1, conv_jydz)
		jydz_loader := db.NewLoader(path.FileInfo(), "jydzb", jydz_reader)
		jydz_loader.Load()

		fmt.Println("导入分工明细表")
		jyfg_reader, _ := f.NewReader("分工表", "A,C,D", 1)
		jyfg_loader := db.NewLoader(path.FileInfo(), "fgmxb", jyfg_reader)
		jyfg_loader.Load()
	}
}

// update_bbmx 更新版本明细
func update_bbmx(db *sqlite.DB) {
	if path := utils.NewPath(config.Home).Find("柜面版本明细*.xlsx"); path != nil {
		fmt.Println("处理文件：", path.Name())
		load_bbmx(db, path)
		export_bbmx(db, path)
	} else {
		fmt.Println("未发现文件：柜面版本明细*.xlsx")
	}
}

// export_bbmx 导出版本明细表
func export_bbmx(db *sqlite.DB, path *utils.Path) {
	file := excel.NewWriter()
	defer file.SaveAs(path.String())
	export_tongji(db, file)
	export_ystm(db, file)
	export_jydzb(db, file)
	export_fgb(db, file)
}

// export_ystm 导出需求条目表
func export_ystm(db *sqlite.DB, w *excel.Writer) {
	fmt.Println("导出需求条目明细表")
	sheet := w.GetSheet("需求条目列表")

	sheet.SetWidth(map[string]float64{
		"A":   20,
		"B:C": 50,
		"D":   20,
		"E":   8,
		"F":   10,
		"G:K": 18,
	})
	sheet.SetColStyle(map[string]string{
		"A:K": "Normal",
	})
	var tcrq string
	if err := db.QueryRow("select max(tcrq)from ystmb").Scan(&tcrq); err != nil {
		return
	}
	query := `select bh,mc,gs,glxt,jsfzr,zt,wbce,hxzc,cszt,bz,ywry from ystmb where tcrq=? order by bh`
	header := "条目编号,功能名称,功能概述,关联系统,负责人,状态,测试人员,核心支持人员,测试状态,情况备注,业务人员"
	rows, err := db.Query(query, tcrq)
	if err != nil {
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.AddTable("A1", header, ch)
}

// export_jydzb 导出新旧交易对照表
func export_jydzb(db *sqlite.DB, w *excel.Writer) {
	fmt.Println("导出新旧交易对照表")
	sheet := w.GetSheet("新旧交易对照表")
	sheet.SetWidth(map[string]float64{
		"A":   20,
		"B":   110,
		"C:D": 12,
		"E":   15,
	})
	sheet.SetColStyle(map[string]string{
		"A:E": "Normal",
	})
	query := `select a.bh,b.mc,a.jym,a.yjym,b.cszt from jydzb a left join ystmb b on a.bh=b.bh order by a.bh,a.yjym`
	header := "条目编号,功能名称,新交易,老交易,进度"
	rows, err := db.Query(query)
	if err != nil {
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.AddTable("A1", header, ch)
}

// export_fgb 导出分工表
func export_fgb(db *sqlite.DB, w *excel.Writer) {
	fmt.Println("导出分工表")
	sheet := w.GetSheet("分工表")
	sheet.SetWidth(map[string]float64{
		"A":   20,
		"B":   110,
		"C:D": 30,
		"E":   15,
	})
	sheet.SetColStyle(map[string]string{
		"A:E": "Normal",
	})
	var tcrq string
	if err := db.QueryRow("select max(tcrq)from ystmb").Scan(&tcrq); err != nil {
		return
	}
	query := `select b.bh,b.mc,a.ywxz,a.jsxz,b.cszt
	from ystmb b left join fgmxb a  on a.bh=b.bh
	where b.tcrq=?
	order by b.bh`
	header := "条目编号,功能名称,业务小组,技术小组,进度"
	rows, err := db.Query(query, tcrq)
	if err != nil {
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.AddTable("A1", header, ch)
}

//go:embed tables/kaifatongji.toml
var kaifatongji string

//go:embed tables/ywxztj.toml
var ywxztj string

func export_tongji(db *sqlite.DB, w *excel.Writer) {
	sheet := w.GetSheet(0)
	sheet.Rename("统计表")
	query := `select b.jsxz,sum(iif(a.cszt like '0%',1,0)),
sum(iif(a.cszt like '1%',1,0)),
sum(iif(a.cszt like '2%',1,0)),
sum(iif(a.cszt like '3%',1,0)),
sum(iif(a.cszt like '4%',1,0)),
count(a.bh) as sl
from ystmb a left join fgmxb b on a.bh=b.bh
where tcrq=?
group by b.jsxz
order by sl desc
`
	var tcrq string
	if err := db.QueryRow("select max(tcrq)from ystmb").Scan(&tcrq); err != nil {
		return
	}
	sheet.SetWidth(map[string]float64{
		"A":   20,
		"B:I": 10,
	})
	rows, err := db.Query(query, tcrq)
	if err != nil {
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	utils.PrintErr(sheet.AddTableToml("A1", kaifatongji, ch))

	query = `select b.ywxz,sum(iif(a.cszt like '0%',1,0)),
sum(iif(a.cszt like '1%',1,0)),
sum(iif(a.cszt like '2%',1,0)),
sum(iif(a.cszt like '3%',1,0)),
sum(iif(a.cszt like '4%',1,0)),
count(a.bh) as sl
from ystmb a left join fgmxb b on a.bh=b.bh
where tcrq=?
group by b.ywxz
order by sl desc
`
	var count int
	db.QueryRow("select count(distinct b.jsxz)from ystmb a left join fgmxb b on a.bh=b.bh where tcrq=?", tcrq).Scan(&count)
	rows, err = db.Query(query, tcrq)
	if err != nil {
		return
	}
	ch = make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	cell, _ := excel.Cell("A", count+5)
	fmt.Println(cell)
	utils.PrintErr(sheet.AddTableToml(cell, ywxztj, ch))
}
