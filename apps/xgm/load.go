package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"peach/data"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"slices"
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
		//loader.Check = false
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
		//loader.Check = false
		loader.Load()
	} else {
		fmt.Println(err)
	}
}

// LoadKfjh 开发计划
func LoadKfjh(db *sqlite.DB, fileinfo fs.FileInfo, book *excel.ExcelBook, ver string) {
	fmt.Println("导入开发计划表")
	if r, err := book.NewReader("开发计划", "A,H:S", 1); err == nil {
		loader := db.NewLoader(fileinfo, "kfjh", r)
		loader.Ver = ver
		//loader.Check = false
		loader.Load()
	} else {
		fmt.Println(err)
	}
}

// Update 更新计划表
func Update(db *sqlite.DB) (err error) {
	path := utils.NewPath("~/Downloads").Find("*新柜面存量交易迁移*.xlsx")
	if path == nil {
		return fmt.Errorf("未找到文件：新柜面存量交易迁移*.xlsx")
	}
	fmt.Println("处理文件：", path.Name())
	ver := utils.Extract(`\d{8}`, path.String())
	if f, err := excel.Open(path); err == nil {
		defer f.Close()
		book := f.ExcelBook
		fileinfo := path.FileInfo()
		LoadKfjh(db, fileinfo, &book, ver)
		LoadXjdzb(db, fileinfo, &book, ver)
		LoadXmjh2(db, fileinfo, &book, ver)
	}
	Update_ytc(db)
	return Export(db, path)
}

// LoadXmjh2 项目计划
func LoadXmjh2(db *sqlite.DB, fileinfo fs.FileInfo, book *excel.ExcelBook, ver string) {
	fmt.Println("导入项目计划表")
	names := make([]string, 0)
	Sheets := []string{"完成表", "计划表", "全量表"}
	for _, name := range book.GetSheetList() {
		if slices.Contains(Sheets, name) {
			names = append(names, name)
		}
	}
	fmt.Println("导入表格：", names)
	if r, err := book.NewReader(names, "A:Q", 1, data.HashFilter(-1, -10, -9, -8, -7, -6, -5, -4, -3, -2)); err == nil {
		loader := db.NewLoader(fileinfo, "xmjh", r)
		loader.Ver = ver
		loader.Method = "insert or replace"
		loader.Clear = false
		//loader.Check = false
		//loader.Test(db)
		loader.Load()

	} else {
		fmt.Println(err)
	}
}

// Load 导入数据文件
func Load(db *sqlite.DB) (err error) {
	path := utils.NewPath("~/Downloads").Find("*新柜面存量交易迁移*.xlsx")
	if path == nil {
		return fmt.Errorf("未找到文件：新柜面存量交易迁移*.xlsx")
	}
	fmt.Println("处理文件：", path.Name())
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
	return
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
		fmt.Println("处理文件：", path.Name())
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
	loader.Ver = strings.Join([]string{ver[:4], ver[4:6], ver[6:]}, "-")
	loader.Load()
}

// Restore 从备份文件中恢复数据
func Restore(db *sqlite.DB) {
	defer utils.TimeIt(time.Now())
	path := utils.NewPath("~/Downloads").Find("新柜面简报*.tgz")
	if path != nil {
		fmt.Println("处理文件：", path.Name())
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

func load_jyjh(db *sqlite.DB) (err error) {
	path := utils.NewPath("~/Downloads").Find("*新柜面剩余交易的迁移计划-*.xlsx")
	if path == nil {
		return fmt.Errorf("未找到 新柜面剩余交易的迁移计划 文件")
	}
	fmt.Println("处理文件：", path.Base())
	f, err := excel.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	r, err := f.NewReader("2026年交易开发计划", "B,R,U,K,AE,AF,AP,AQ,AR,AH,AI,AI,AJ", 1)
	if err != nil {
		return
	}
	d := data.NewData()
	go r.Read(d)
	d.Println()
	return
}
