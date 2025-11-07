package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"
	"text/template"
)

// conv_gzb 转换故障表的数据
func conv_gzb(src []string) (dest []string, err error) {
	for _, k := range []int{4, 10} {
		src[k] = excel.FormatDate(src[k], "2006-01-02")
	}
	dest = src
	return
}

// LoadWtgzb 导入问题跟踪表数据
func LoadWtgzb(db *sqlite.DB, file utils.File) (err error) {
	var (
		f    io.ReadCloser
		book *excel.ExcelBook
	)
	if f, err = file.Open(); err == nil {
		defer f.Close()
		if book, err = excel.NewExcelBook(f, file.FileInfo().Name()); err == nil {
			load_wtgzb(db, book, file.FileInfo())
		}
	}
	return
}

// load_wtgzb 导入问题跟踪表
func load_wtgzb(db *sqlite.DB, book *excel.ExcelBook, fileinfo fs.FileInfo) (err error) {
	name := fileinfo.Name()
	ver := utils.Extract(`\d{8}`, name)
	fmt.Println("处理文件：", name, "Version:", ver)
	if r, err := book.NewReader(0, "A:M", 1, conv_gzb); err == nil {
		loader := db.NewLoader(fileinfo, "wtgzb", r)
		loader.Ver = strings.Join([]string{ver[:4], ver[4:6], ver[6:]}, "-")
		return loader.Load()
	} else {
		return err
	}
}

type wtTongji struct {
	Zongshu      int
	Yanzhongxing map[string]int
	Fenlei       map[string]int
	Zhuangtai    map[string]int
}

type wtReport struct {
	Baogaoqi string
	Dangqi   wtTongji
	Heji     wtTongji
}

// report_wenti 统计分支行报送的问题情况
func report_wenti(db *sqlite.DB, bgq string) (err error) {
	var (
		rows  *sqlite.Rows
		key   string
		value int
	)
	wt := &wtReport{Baogaoqi: bgq}
	if err = db.QueryRow("select count(*)from wtgzb").Scan(&wt.Heji.Zongshu); err != nil {
		return
	}
	if rows, err = db.Query("select yzx,count(xh)from wtgzb group by yzx"); err != nil {
		return
	}
	wt.Heji.Yanzhongxing = make(map[string]int)
	for rows.Next() {
		rows.Scan(&key, &value)
		wt.Heji.Yanzhongxing[key] = value
	}

	if rows, err = db.Query("select zt,count(xh)from wtgzb group by zt"); err != nil {
		return
	}
	wt.Heji.Zhuangtai = make(map[string]int)
	for rows.Next() {
		rows.Scan(&key, &value)
		wt.Heji.Zhuangtai[key] = value
	}

	if rows, err = db.Query("select wtfl,count(xh)from wtgzb group by wtfl"); err != nil {
		return
	}
	wt.Heji.Fenlei = make(map[string]int)
	for rows.Next() {
		rows.Scan(&key, &value)
		wt.Heji.Fenlei[key] = value
	}
	//utils.PrintStruct(wt)

	query := "where 1=1"
	switch bgq {
	case "本月":
		query = "where strftime('%Y-%m',tcrq)=strftime('%Y-%m',date('now')) "
	case "上月":
		query = "where strftime('%Y-%m',tcrq)=strftime('%Y-%m',date('now','-1 months')) "
	case "本周":
		query = "where strftime('%W',tcrq)=strftime('%W',date('now')) "
	case "上周":
		query = "where strftime('%W',tcrq)=strftime('%W',date('now','-7 days')) "
	}
	if err = db.QueryRow(fmt.Sprintf("select count(*)from wtgzb %s", query)).Scan(&wt.Dangqi.Zongshu); err != nil {
		return
	}
	if rows, err = db.Query(fmt.Sprintf("select yzx,count(xh)from wtgzb %s group by yzx", query)); err != nil {
		return
	}
	wt.Dangqi.Yanzhongxing = make(map[string]int)
	for rows.Next() {
		rows.Scan(&key, &value)
		wt.Dangqi.Yanzhongxing[key] = value
	}
	if rows, err = db.Query(fmt.Sprintf("select zt,count(xh)from wtgzb %s group by zt", query)); err != nil {
		return
	}
	wt.Dangqi.Zhuangtai = make(map[string]int)
	for rows.Next() {
		rows.Scan(&key, &value)
		wt.Dangqi.Zhuangtai[key] = value
	}

	if rows, err = db.Query(fmt.Sprintf("select wtfl,count(xh)from wtgzb %s group by wtfl", query)); err != nil {
		return
	}
	wt.Dangqi.Fenlei = make(map[string]int)
	for rows.Next() {
		rows.Scan(&key, &value)
		wt.Dangqi.Fenlei[key] = value
	}

	tmp, _ := template.ParseFS(templateFS, "template/wenti.txt")
	tmp.Execute(os.Stdout, wt)
	return
}
