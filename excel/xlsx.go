package excel

import (
	"errors"
	"io"
	"iter"
	"peach/data"
	"peach/utils"
	"strings"

	"github.com/xuri/excelize/v2"
)

// XlsxBook 对 xlsx 工作簿的封装
type xlsxBook struct {
	*excelize.File
}

type ErrSheetNotExist excelize.ErrSheetNotExist

// UseCols 设置选取哪些列，表达式为：A:C,E:F,H
func UseCols(cols string) func([]string) ([]string, error) {
	columns := make([]int, 0)
	for c := range strings.SplitSeq(cols, ",") {
		if strings.Contains(c, ":") {
			x := strings.Split(c, ":")
			a, err := excelize.ColumnNameToNumber(x[0])
			utils.CheckFatal(err)
			b, err := excelize.ColumnNameToNumber(x[1])
			utils.CheckFatal(err)
			for i := a; i <= b; i++ {
				columns = append(columns, i-1)
			}
		} else {
			i, err := excelize.ColumnNameToNumber(c)
			utils.CheckFatal(err)
			columns = append(columns, i-1)
		}
	}
	return data.Include(columns...)
}

// NewXlsxBook  XlsxBook 的构造函数
func newXlsxBook(reader io.Reader, opts ...excelize.Options) (r *xlsxBook, err error) {
	book, err := excelize.OpenReader(reader, opts...)
	if err != nil {
		return
	}
	r = &xlsxBook{book}
	return
}

func (b *xlsxBook) IterRows(sheetIdx int, skipRows int) iter.Seq[[]string] {
	return func(yield func([]string) bool) {
		name := b.GetSheetName(sheetIdx)
		rows, _ := b.Rows(name)
		for range skipRows {
			rows.Next()
		}
		for rows.Next() {
			columns, err := rows.Columns()
			if err != nil || !yield(columns) {
				break
			}
		}
	}
}

// ReedSheet 逐行读取数据
func (b *xlsxBook) ReadSheet(num int, skipRows int, ch chan<- []any, cvfns ...data.ConvertFunc) {
	defer close(ch)
	name := b.GetSheetName(num)
	rows, _ := b.Rows(name)
	for range skipRows {
		rows.Next()
	}
	for rows.Next() {
		columns, err := rows.Columns()
		if err == nil {
			result, err := data.Convert(cvfns...)(columns)
			if result != nil && err == nil {
				ch <- result
			}
		}
	}
}

// GetValues 读取所有数据
func (b *xlsxBook) GetValues(num int) (values [][]string, err error) {
	name := b.GetSheetName(num)
	if name == "" {
		err = errors.New("sheet 不存在")
		return
	}
	return b.GetRows(name)
}
