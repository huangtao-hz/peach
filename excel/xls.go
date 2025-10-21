package excel

import (
	"errors"
	"io"
	"iter"
	"peach/data"

	"github.com/extrame/xls"
)

// XlsBook 对 xls 工作簿的封装
type xlsBook struct {
	*xls.WorkBook
}

// newXlsBook XlsBook 的构造函数
func newXlsBook(reader io.ReadSeeker) (r *xlsBook, err error) {
	book, err := xls.OpenReader(reader, "utf8")
	if err != nil {
		return
	}
	r = &xlsBook{book}
	return
}

// GetSheetList 获取工作表列表
func (b *xlsBook) GetSheetList() (names []string) {
	num := b.NumSheets()
	names = make([]string, num)
	for i := range num {
		names[i] = b.GetSheet(i).Name
	}
	return
}

// GetRows 获取每一行的数据
func (b *xlsBook) IterRows(sheetIdx int, skipRows int) iter.Seq[[]string] {
	return func(yield func([]string) bool) {
		sheet := b.GetSheet(sheetIdx)
		for i := skipRows; i <= int(sheet.MaxRow); i++ {
			row := sheet.Row(i)
			line := make([]string, row.LastCol()-row.FirstCol())
			for i := row.FirstCol(); i < row.LastCol(); i++ {
				line[i] = row.Col(i)
			}
			if !yield(line) {
				break
			}
		}
	}
}

// ReadSheet 读取 sheet 数据
func (b *xlsBook) ReadSheet(num int, skipRows int, ch chan<- []any, cvfns ...data.ConvertFunc) {
	defer close(ch)
	sheet := b.GetSheet(num)
	rowcount := int(sheet.MaxRow) + 1
	for i := skipRows; i < rowcount; i++ {
		row := sheet.Row(i)
		line := make([]string, row.LastCol()-row.FirstCol())
		for i := row.FirstCol(); i < row.LastCol(); i++ {
			line[i] = row.Col(i)
		}
		result, err := data.Convert(cvfns...)(line)
		if row != nil && err == nil {
			ch <- result
		}
	}
}

// GetValues 获取数据
func (b *xlsBook) GetValues(num int) (data [][]string, err error) {
	var rowcount int
	if num < 0 && num >= b.NumSheets() {
		err = errors.New("表格录入错误")
		return
	}
	sheet := b.GetSheet(num)
	rowcount = int(sheet.MaxRow) + 1
	data = make([][]string, rowcount)
	for r := range rowcount {
		row := sheet.Row(r)
		line := make([]string, row.LastCol()-row.FirstCol())
		for i := row.FirstCol(); i < row.LastCol(); i++ {
			line[i] = row.Col(i)
		}
		data[r] = line
	}
	return
}
