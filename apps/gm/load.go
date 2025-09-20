package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"peach/data"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"
)

// 打印当前数据版本
func PrintVersion(db *sqlite.DB) {
	var ver string
	if err := db.QueryRow("select ver from loadfile where name='xmjh'").Scan(&ver); err == nil {
		fmt.Println("当前数据版本：", ver)
	}
}

// 新旧交易对照表
func LoadXjdzb(db *sqlite.DB, path string, r *excel.ExcelReader, ver string) {
	ch := make(chan []any, 100)
	loader := sqlite.LoadFile(path, "jydzb", ch)
	loader.Ver = ver
	go r.ReadSheet(7, 1, ch, data.FixedColumn(7))
	loader.Load(db)
	//utils.ChPrintln(ch)
}

// 项目计划
func LoadXmjh(db *sqlite.DB, path string, r *excel.ExcelReader, ver string) {
	ch := make(chan []any, 100)
	loader := sqlite.LoadFile(path, "xmjh", ch)
	loader.Ver = ver
	go r.ReadSheet(6, 1, ch, data.FixedColumn(16))
	loader.Load(db)
	//utils.ChPrintln(ch)
}

func conv_kfjh(s []string) (d []string, err error) {
	idx := []int{0, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	for len(s) < 24 {
		s = append(s, "")
	}
	d = make([]string, len(idx))
	for i := range len(idx) {
		d[i] = s[idx[i]]
	}
	return
}

// 开发计划
func LoadKfjh(db *sqlite.DB, path string, r *excel.ExcelReader, ver string) {
	ch := make(chan []any, 100)
	loader := sqlite.LoadFile(path, "kfjh", ch)
	loader.Ver = ver
	go r.ReadSheet(3, 1, ch, conv_kfjh)
	loader.Load(db)
}

// 导入数据文件
func Load(db *sqlite.DB) {
	path := utils.NewPath("~/Downloads").Find("*新柜面存量交易迁移*.xlsx")
	if path != "" {
		fmt.Println("处理文件：", utils.NewPath(path).FileInfo().Name())
	}
	ver := utils.Extract(`\d{8}`, path)
	fmt.Println("Version:", ver)
	r, err := excel.NewExcelReader(path)
	utils.CheckFatal(err)
	defer r.Close()
	LoadKfjh(db, path, r, ver)
	LoadXjdzb(db, path, r, ver)
	LoadXmjh(db, path, r, ver)
	load_gzb(db)
}

// conv_gzb 转换故障表的数据
func conv_gzb(src []string) (dest []string, err error) {
	for _, k := range []int{4, 10} {
		src[k] = excel.FormatDate(src[k])
	}
	dest = src
	return
}

// load_gzb 导入问题跟踪表数据
func load_gzb(db *sqlite.DB) {
	path := utils.NewPath("~/Downloads").Find("*数智综合运营系统问题跟踪表*.xlsx")
	if path != "" {
		fmt.Println("处理文件：", utils.NewPath(path).FileInfo().Name())
	} else {
		fmt.Println("未找到文件！")
		return
	}
	ver := utils.Extract(`\d{8}`, path)
	fmt.Println("Version:", ver)
	r, err := excel.NewExcelReader(path)
	utils.CheckFatal(err)
	defer r.Close()

	ch := make(chan []any, 100)
	go r.ReadSheet(0, 1, ch, data.FixedColumn(13), conv_gzb)
	utils.ChPrintln(ch)
}

// Restore 从备份文件中恢复数据
func Restore(db *sqlite.DB) {
	path := utils.NewPath("~/Downloads").Find("新柜面简报*.tgz")
	if path != "" {
		fmt.Println("处理文件：", utils.NewPath(path).FileInfo().Name())
	} else {
		fmt.Println("未找到文件！")
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return
	}
	t := tar.NewReader(r)
	for header, err := t.Next(); err != io.EOF; header, err = t.Next() {
		name := header.FileInfo().Name()
		if strings.Contains(name, "新柜面存量交易迁移计划") {
			fmt.Println("处理文件：", name)
			ver := utils.Extract(`\d{8}`, name)
			fmt.Println("Version:", ver)
			r, err := excel.NewXlsxFile(t)
			utils.CheckFatal(err)
			fmt.Println(r.GetSheetList())
			ch := make(chan []any, 100)
			go r.ReadSheet(4, 1, ch, data.FixedColumn(16))
			utils.ChPrintln(ch)
		}
	}

}
