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
	jyfg_loader.Method = "insert or replace"
	jyfg_loader.Clear = false
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
	export_tongji(db, file)
	for path := range strings.SplitSeq("bb_ystm,bb_xjdzb,bb_fgb,bb_xmryb", ",") {
		path = fmt.Sprintf("%s.toml", path)
		utils.CheckFatal(ExportReport(db, file, path))
	}
	file.SaveAs(path.String())
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
