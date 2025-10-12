package excel

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/huangtao-hz/excelize"
)

var (
	dateFormat1 *regexp.Regexp = regexp.MustCompile(`^\d{2}-\d{2}-\d{2}$`)
	dateFormat2 *regexp.Regexp = regexp.MustCompile(`^\d+$`)
	dateFormat3 *regexp.Regexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
)

// FormatDate 格式化 Excel 文件中的日期
// 将日期统格式化成 YYYY-MM-DD 格式，
// 无法格式化的，沿用原来的值
func FormatDate(s string, layout string) string {
	if d, err := AToDate(s, layout); err == nil {
		return d
	}
	return s
}

// AToDate 格式化日期
func AToDate(s string, layout string) (d string, err error) {
	var (
		t time.Time
		f float64
	)
	if dateFormat1.MatchString(s) {
		t, err = time.Parse("01-02-06", s)
	} else if dateFormat2.MatchString(s) {
		if f, err = strconv.ParseFloat(s, 64); err == nil {
			t, err = excelize.ExcelDateToTime(f, false)
		}
	} else if dateFormat3.MatchString(s) {
		t, err = time.Parse("2006-01-02", s)
	} else {
		err = fmt.Errorf("无法将 %v 转换成日期", s)
	}
	if err == nil {
		d = t.Format(layout)
	}
	return
}

// Cell 把坐标转换成单位元格
func Cell(col any, row int) (cell string, err error) {
	switch column := col.(type) {
	case int:
		cell, err = excelize.CoordinatesToCellName(column, row)
	case string:
		cell, err = excelize.JoinCellName(column, row)
	default:
		err = fmt.Errorf("Cell(%v,%d)不是有效的单元格格式", col, row)
	}
	return
}

// RangeToCoordinates 将区间转换为坐标
func RangeToCoordinates(rng string) (col1, row1, col2, row2 int, err error) {
	cells := strings.Split(rng, ":")
	if len(cells) > 2 {
		err = fmt.Errorf("%s 非有效的单元格格式", rng)
	} else if len(cells) == 1 {
		cells = append(cells, cells[0])
	}
	col1, row1, err = excelize.CellNameToCoordinates(cells[0])
	if err != nil {
		return
	}
	col2, row2, err = excelize.CellNameToCoordinates(cells[1])
	if err != nil {
		return
	}
	return
}

// Range2Cells 将 range 转换为 Cells
func Range2Cells(rng string) (cell1, cell2 string, err error) {
	cells := strings.Split(rng, ":")
	count := len(cells)
	if count < 1 || count > 2 {
		err = fmt.Errorf("%s 不是有效的单元格", rng)
	} else if count == 1 {
		cell1, cell2 = cells[0], cells[0]
	} else {
		cell1, cell2 = cells[0], cells[1]
	}
	return
}
