package main

import (
	"fmt"
	"peach/sqlite"
	"peach/utils"
	"strings"
)

// convdate 转换日期
func convdate(d *string) {
	if strings.Contains(*d, "/") {
		*d = strings.ReplaceAll(*d, "/", "-")
	}
	if len(*d) > 0 && len(*d) < 10 {
		*d = ""
	}
}

// conv_kfjh 转换开发计划
func conv_kfjh(src []string) (dest []string, err error) {
	if src[0] == "" {
		return
	} else if len(src[0]) < 4 {
		src[0] = fmt.Sprintf("%04s", src[0])
	}
	for i := 1; i < 10; i++ {
		src[i] = strings.TrimSpace(src[i])
	}
	for i := 10; i < 14; i++ {
		convdate(&src[i])
	}
	dest = src
	return
}

// load_kfjh 导入开发计划表
func load_kfjh(db *sqlite.DB) (err error) {
	if path := utils.NewPath(config.Home).Find("*开发计划*.xlsx"); path != nil {
		fmt.Println("处理文件：", path.Base())
		date := utils.Extract(`\d{8}`, path.Base())
		date = strings.Join([]string{date[:4], date[4:6], date[6:]}, "-")
		return db.LoadExcelFile(loaderFS, "loader/kfjh.toml", path, date, conv_kfjh)
	}
	return fmt.Errorf("未找到 开发计划 文件")
}
