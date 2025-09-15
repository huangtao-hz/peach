package excel

import (
	"errors"
	"io"
	"peach/data"

	"github.com/extrame/xls"
)

type XlsFile struct {
	*xls.WorkBook
}

func NewXlsFile(reader io.ReadSeeker) (r *XlsFile, err error) {
	book, err := xls.OpenReader(reader, "utf8")
	if err != nil {
		return
	}
	r = &XlsFile{book}
	return
}

func (r *XlsFile) GetSheetList() (names []string) {
	num := r.NumSheets()
	names = make([]string, num)
	for i := range num {
		names[i] = r.GetSheet(i).Name
	}
	return
}

func (rd *XlsFile) ReadSheet(num int, skipRows int, ch chan<- []any, cvfns ...data.ConvertFunc) {
	defer close(ch)
	sheet := rd.GetSheet(num)
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

func (rd *XlsFile) GetValues(num int) (data [][]string, err error) {
	var rowcount int
	if num < 0 && num >= rd.NumSheets() {
		err = errors.New("表格录入错误")
		return
	}
	sheet := rd.GetSheet(num)
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
