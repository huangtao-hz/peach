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

// get_tcrq 获取最近的投产日期
func get_tcrq(db *sqlite.DB) (tcrq string, err error) {
	err = db.QueryRow("select max(tcrq)from ystmb").Scan(&tcrq)
	return
}

// load_bbmx 导入版本明细数据
func load_bbmx(db *sqlite.DB, path utils.File) {
	name := path.FileInfo().Name()
	r, err := path.Open()
	if err != nil {
		return
	}
	defer r.Close()
	f, err := excel.NewExcelBook(r, name)
	if err != nil {
		return
	}
	fmt.Println("导入需求条目表")
	date := utils.Extract(`\d{8}`, name)
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

	fmt.Println("导入项目人员表")
	xmry_reader, err := f.NewReader("项目人员表", "A:C", 1)
	if err != nil {
		return
	}
	xmry_loader := db.NewLoader(path.FileInfo(), "xmryb", xmry_reader)
	xmry_loader.Load()
}

// update_bbmx 更新版本明细
func update_bbmx(db *sqlite.DB) {
	if path := utils.NewPath(config.Home).Find("*版本条目明细*.xlsx"); path != nil {
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
	export_xmryb(db, file)
}

// export_ystm 导出需求条目表
func export_ystm(db *sqlite.DB, w *excel.Writer) {
	var err error
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

	tcrq, err := get_tcrq(db)
	if err != nil {
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

// export_xmryb 导出项目人员表
func export_xmryb(db *sqlite.DB, w *excel.Writer) {
	fmt.Println("导出项目人员表")
	sheet := w.GetSheet("项目人员表")
	sheet.SetWidth(map[string]float64{
		"A": 20,
		"B": 30,
		"C": 30,
	})
	sheet.SetColStyle(map[string]string{
		"A:C": "Normal",
	})

	query := `select * from xmryb order by lb,xz`
	header := "小组,类别,姓名"
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

	tcrq, err := get_tcrq(db)
	if err != nil {
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

var (
	//go:embed tables/kaifatongji.toml
	kaifatongji string
	//go:embed tables/ywxztj.toml
	ywxztj string
	//go:embed query/kaifatongji.sql
	kaifa_query string
	//go:embed query/ywtongji.sql
	yw_query string
)

func export_tongji(db *sqlite.DB, w *excel.Writer) {
	sheet := w.GetSheet("统计表")
	tcrq, err := get_tcrq(db)
	if err != nil {
		return
	}
	sheet.SetWidth(map[string]float64{
		"A":   20,
		"B:Z": 10,
	})
	rows, err := db.Query(kaifa_query, tcrq)
	if err != nil {
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	utils.PrintErr(sheet.AddTableToml("A1", kaifatongji, ch))

	var count int
	db.QueryRow(fmt.Sprintf("select count(*)from (%s)", kaifa_query), tcrq).Scan(&count)
	rows, err = db.Query(yw_query, tcrq)
	if err != nil {
		return
	}
	ch = make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.SkipRows(2)
	utils.PrintErr(sheet.AddTableToml("A", ywxztj, ch))
}
