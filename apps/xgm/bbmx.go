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
	jydz_loader.Method = "insert or replace"
	jydz_loader.Clear = false
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
		utils.CheckFatal(ExportXlsx(db, path.String(), "bb_kjtj,bb_ywtj,bb_ystm,bb_xjdzb,bb_fgb,bb_xmryb"))
	} else {
		fmt.Println("未发现文件：柜面版本明细*.xlsx")
	}
}
