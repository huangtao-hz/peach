package main

import (
	"fmt"
	"io"
	"peach/data"
	"peach/sqlite"
	"peach/utils"
	"slices"
	"strings"

	"github.com/huangtao-hz/excelize"
)

func conv_jym(jym *string) {
	if utils.FullMatch(`\d{3}`, *jym) {
		*jym = fmt.Sprintf("%04s", *jym)
	}
}

// excelReader 读取 Excel文件
type excelReader struct {
	name string
	file utils.File
}

// GetMergedCells
func GetMergedCells(file *excelize.File, sheet string) (cells map[string]string, err error) {
	var (
		c1, r1, c2, r2 int
		mcs            []excelize.MergeCell
	)
	cells = make(map[string]string)
	if mcs, err = file.GetMergeCells(sheet, false); err == nil {
		for _, mc := range mcs {
			c1, r1, _ = excelize.CellNameToCoordinates(mc.GetStartAxis())
			c2, r2, _ = excelize.CellNameToCoordinates(mc.GetEndAxis())
			for r := r1; r <= r2; r++ {
				for c := c1; c <= c2; c++ {
					cell, _ := excelize.CoordinatesToCellName(c, r)
					cells[cell] = mc.GetCellValue()
				}
			}
		}
	}
	return
}

// read_excel 读取 excel 文件
func (er *excelReader) Read(d *data.Data) {
	defer close(d.Data)
	name, file := er.name, er.file
	rq := utils.Extract(`\d{8}`, name)
	var (
		r    io.ReadCloser
		f    *excelize.File
		err  error
		rows *excelize.Rows
	)
	if r, err = file.Open(); err != nil {
		return
	}
	defer r.Close()
	opt := excelize.Options{ShortDatePattern: "yyyy-mm-dd"}
	if f, err = excelize.OpenReader(r, opt); err != nil {
		return
	}
	sheet := f.GetSheetName(0)
	mgcells, _ := GetMergedCells(f, sheet)
	rq2 := utils.Extract(`\d{4}-?\d{2}-?\d{2}`, sheet)
	rq2 = strings.ReplaceAll(rq2, "-", "")
	if rq != rq2 {
		fmt.Println("Error:", rq, rq2, sheet)
	}
	rq = strings.Join([]string{rq[:4], rq[4:6], rq[6:]}, "-")
	fmt.Println("日期：", rq)

	if rows, err = f.Rows(sheet); err != nil {
		return
	}

	for row := 1; rows.Next(); row++ {
		if row > 1 {
			cols, _ := rows.Columns()
			if len(cols) == 0 {
				continue
			} else if len(cols) < 11 {
				cols = append(cols, slices.Repeat([]string{""}, 11-len(cols))...)
			}
			for c, v := range cols {
				if v == "" && c < 11 {
					cell, _ := excelize.CoordinatesToCellName(c, row)
					cols[c], _ = mgcells[cell]
				}
			}
			cols[0] = rq
			conv_jym(&cols[2])
			d.Data <- utils.Slice(cols)
		}
	}
}

// load_bbsm 导入单个文件
func load_bbsm(db *sqlite.DB, name string, file utils.File) (err error) {
	//fmt.Println("处理文件：", name)
	r := &excelReader{name: name, file: file}
	rq := utils.Extract(`\d{8}`, name)
	loader := db.NewLoader(file.FileInfo(), "bbsm", r)
	loader.Check = true
	loader.Clear = true
	loader.Method = "insert"
	rq = strings.Join([]string{rq[:4], rq[4:6], rq[6:]}, "-")
	loader.ClearSQL = fmt.Sprintf("delete from bbsm where rq='%s'", rq)
	return loader.Load()
}

// load_all 导入版本说明压缩包
func load_all(db *sqlite.DB) {
	if path := utils.NewPath("~/Downloads").Find("投产版本说明*.zip"); path != nil {
		for name, file := range path.IterZip() {
			utils.PrintErr(load_bbsm(db, name, file))
		}
	}
}
