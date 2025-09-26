package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"
	"time"
)

// PrintVersion 打印当前数据版本
func PrintVersion(db *sqlite.DB) {
	var ver string
	if err := db.QueryRow("select ver from loadfile where name='xmjh'").Scan(&ver); err == nil {
		fmt.Println("当前数据版本：", ver)
	}
}

// LoadXjdzb 新旧交易对照表
func LoadXjdzb(db *sqlite.DB, fileinfo fs.FileInfo, book *excel.ExcelBook, ver string) {
	fmt.Println("导入新旧交易对照表")
	if r, err := book.NewReader("投产交易一览表", "A:G", 1); err == nil {
		loader := db.NewLoader(fileinfo, "xjdz", r)
		loader.Ver = ver
		loader.Check = false
		loader.Load()
	} else {
		fmt.Println(err)
	}
}

// LoadXmjh 项目计划
func LoadXmjh(db *sqlite.DB, fileinfo fs.FileInfo, book *excel.ExcelBook, ver string) {
	fmt.Println("导入项目计划表")
	if r, err := book.NewReader("全量表", "A:P", 1); err == nil {
		loader := db.NewLoader(fileinfo, "xmjh", r)
		loader.Ver = ver
		loader.Check = false
		loader.Load()
	} else {
		fmt.Println(err)
	}
}

// LoadKfjh 开发计划
func LoadKfjh(db *sqlite.DB, fileinfo fs.FileInfo, book *excel.ExcelBook, ver string) {
	fmt.Println("导入开发计划表")
	if r, err := book.NewReader("开发计划", "A,M:X", 1); err == nil {
		loader := db.NewLoader(fileinfo, "kfjh", r)
		loader.Ver = ver
		loader.Check = false
		loader.Load()
	} else {
		fmt.Println(err)
	}
}

// Load 导入数据文件
func Load(db *sqlite.DB) {
	path := utils.NewPath("~/Downloads").Find("*新柜面存量交易迁移*.xlsx")
	if path != nil {
		fmt.Println("处理文件：", path.FileInfo().Name())
	}
	ver := utils.Extract(`\d{8}`, path.String())
	fmt.Println("Version:", ver)
	f, err := excel.NewExcelFile(path.String())
	utils.CheckFatal(err)
	defer f.Close()
	book := f.ExcelBook
	fileinfo := path.FileInfo()
	LoadKfjh(db, fileinfo, &book, ver)
	LoadXjdzb(db, fileinfo, &book, ver)
	LoadXmjh(db, fileinfo, &book, ver)
	Update_ytc(db)
	//load_gzb(db)
}

// conv_gzb 转换故障表的数据
func conv_gzb(src []string) (dest []string, err error) {
	for _, k := range []int{4, 10} {
		src[k] = excel.FormatDate(src[k], "2006-01-02")
	}
	dest = src
	return
}

// LoadWtgzb 导入问题跟踪表数据
func LoadWtgzb(db *sqlite.DB) {
	path := utils.NewPath("~/Downloads").Find("*数智综合运营系统问题跟踪表*.xlsx")
	if path != nil {
		fmt.Println("处理文件：", path.FileInfo().Name())
	}
	f, err := excel.NewExcelFile(path.String())
	utils.CheckFatal(err)
	defer f.Close()
	load_wtgzb(db, &f.ExcelBook, path.FileInfo())
}

// load_wtgzb 导入问题跟踪表
func load_wtgzb(db *sqlite.DB, book *excel.ExcelBook, fileinfo fs.FileInfo) {
	ver := utils.Extract(`\d{8}`, fileinfo.Name())
	fmt.Println("Version:", ver)
	r, err := book.NewReader(0, "A:M", 1, conv_gzb)
	utils.CheckFatal(err)
	loader := db.NewLoader(fileinfo, "wtgzb", r)
	loader.Ver = ver
	loader.Check = false
	loader.Load()
}

// Restore 从备份文件中恢复数据
func Restore(db *sqlite.DB) {
	defer utils.TimeIt(time.Now())
	path := utils.NewPath("~/Downloads").Find("新柜面简报*.tgz")
	if path != nil {
		fmt.Println("处理文件：", path.FileInfo().Name())
	}
	f, err := path.Open()
	utils.CheckFatal(err)
	defer f.Close()
	r, err := gzip.NewReader(f)
	utils.CheckFatal(err)
	defer r.Close()
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
		} else if strings.Contains(name, "数智综合运营系统问题跟踪表") {
			fmt.Println("处理文件：", name)
			book, err := excel.NewExcelBook(t, name)
			utils.CheckFatal(err)
			load_wtgzb(db, book, fileinfo)
		}
	}
	Update_ytc(db)
}
