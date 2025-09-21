package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"peach/data"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"
)

// PrintVersion 打印当前数据版本
func PrintVersion(db *sqlite.DB) {
	var ver string
	if err := db.QueryRow("select ver from loadfile where name='xmjh'").Scan(&ver); err == nil {
		fmt.Println("当前数据版本：", ver)
	}
}

type ExcelReader interface {
	Read(sheets any, skipRows int, ch chan<- []any, cvfns ...data.ConvertFunc)
}

// LoadXjdzb 新旧交易对照表
func LoadXjdzb(db *sqlite.DB, fileinfo fs.FileInfo, r ExcelReader, ver string) {
	fmt.Println("导入新旧交易对照表")
	ch := make(chan []any, 100)
	loader := sqlite.NewLoader(fileinfo, "jydzb", ch)
	loader.Ver = ver
	go r.Read("投产交易一览表", 1, ch, data.FixedColumn(7))
	loader.Load(db)
	//utils.ChPrintln(ch)
}

// 项目计划
func LoadXmjh(db *sqlite.DB, fileinfo fs.FileInfo, r ExcelReader, ver string) {
	fmt.Println("导入项目计划表")
	ch := make(chan []any, 100)
	loader := sqlite.NewLoader(fileinfo, "xmjh", ch)
	loader.Ver = ver
	//loader.Check = false
	go r.Read("全量表", 1, ch, data.FixedColumn(16))
	loader.Load(db)
	//utils.ChPrintln(ch)
}

// LoadKfjh 开发计划
func LoadKfjh(db *sqlite.DB, fileinfo fs.FileInfo, r ExcelReader, ver string) {
	fmt.Println("导入开发计划表")
	ch := make(chan []any, 100)
	loader := sqlite.NewLoader(fileinfo, "kfjh", ch)
	loader.Ver = ver
	go r.Read("开发计划", 1, ch, excel.UseCols("A,M:X"))
	loader.Load(db)
}

// Load 导入数据文件
func Load(db *sqlite.DB) {
	path := utils.NewPath("~/Downloads").Find("*新柜面存量交易迁移*.xlsx")
	if path != "" {
		fmt.Println("处理文件：", utils.NewPath(path).FileInfo().Name())
	}
	ver := utils.Extract(`\d{8}`, path)
	fmt.Println("Version:", ver)
	r, err := excel.NewExcelFile(path)
	utils.CheckFatal(err)
	defer r.Close()
	fileinfo := utils.NewPath(path).FileInfo()
	LoadKfjh(db, fileinfo, r, ver)
	LoadXjdzb(db, fileinfo, r, ver)
	LoadXmjh(db, fileinfo, r, ver)
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
	r, err := excel.NewExcelFile(path)
	utils.CheckFatal(err)
	defer r.Close()

	ch := make(chan []any, 100)
	go r.Read(0, 1, ch, data.FixedColumn(13), conv_gzb)
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
		fileinfo := header.FileInfo()
		name := fileinfo.Name()
		if strings.Contains(name, "新柜面存量交易迁移计划") {
			fmt.Println("处理文件：", name)
			ver := utils.Extract(`\d{8}`, name)
			fmt.Println("Version:", ver)
			r, err := excel.NewExcelBook(t, name)
			utils.CheckFatal(err)
			LoadKfjh(db, fileinfo, r, ver)
			LoadXjdzb(db, fileinfo, r, ver)
			LoadXmjh(db, fileinfo, r, ver)
		}
	}

}
