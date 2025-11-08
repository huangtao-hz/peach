package main

import (
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
	f, err := excel.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Println("导入需求条目表")
	date := utils.Extract(`\d{8}`, name)
	date = strings.Join([]string{date[:4], date[4:6], date[6:]}, "-")
	conv_ystm := func(src []string) (dest []string, err error) {
		if src[0] == "" {
			return
		}
		dest = []string{date}
		dest = append(dest, src...)
		return
	}
	utils.CheckErr(db.LoadExcel(loaderFS, "loader/bb_ystm.toml", &f.ExcelBook, path.FileInfo(), conv_ystm))

	fmt.Println("导入新旧对照表")
	conv_jydz := func(src []string) (dest []string, err error) {
		if src[0] == "" {
			return
		}
		conv_jym(&src[1])
		dest = src
		return
	}
	utils.CheckErr(db.LoadExcel(loaderFS, "loader/bb_jydzb.toml", &f.ExcelBook, path.FileInfo(), conv_jydz))
	fmt.Println("导入分工明细表")
	utils.CheckErr(db.LoadExcel(loaderFS, "loader/bb_fgb.toml", &f.ExcelBook, path.FileInfo()))
	fmt.Println("导入项目人员表")
	utils.CheckErr(db.LoadExcel(loaderFS, "loader/bb_xmryb.toml", &f.ExcelBook, path.FileInfo()))
}

// update_bbmx 更新版本明细
func update_bbmx(db *sqlite.DB) {
	if path := utils.NewPath(config.Home).Find("*版本条目明细*.xlsx"); path != nil {
		fmt.Println("处理文件：", path.Name())
		load_bbmx(db, path)
		utils.CheckFatal(ExportXlsx(db, path.String(), "bb_kjtj,bb_ywrytj,bb_ywtj,bb_ystm,bb_xjdzb,bb_fgb,bb_xmryb"))
	} else {
		fmt.Println("未发现文件：柜面版本明细*.xlsx")
	}
}
