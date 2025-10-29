package main

import (
	"fmt"
	"peach/excel"
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

// conv_kfzt 开发状态转换
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

func load_kfjh(db *sqlite.DB) (err error) {
	path := utils.NewPath(config.Home).Find("*开发计划*.xlsx")
	if path == nil {
		return fmt.Errorf("未找到 开发计划 文件")
	}
	fmt.Println("处理文件：", path.Base())
	f, err := excel.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := f.NewReader("柜面核心类交易开发计划", "B,R,U,K,H,AE,AF,AP,AW,BD,AG,AH,BM,BS", 1, conv_kfjh)
	if err != nil {
		return err
	}
	loader := db.NewLoader(path.FileInfo(), "kfjh", r)
	loader.Method = "insert or replace"
	loader.Clear = false
	loader.Load()
	return
}
