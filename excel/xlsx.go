package excel

import (
	"errors"
	"io"
	"peach/data"
	"time"

	"github.com/xuri/excelize/v2"
)

type XlsxFile struct {
	*excelize.File
}

type ErrSheetNotExist excelize.ErrSheetNotExist

var ColumnNameToNumber = excelize.ColumnNameToNumber

func NewXlsxFile(reader io.Reader, opts ...excelize.Options) (r *XlsxFile, err error) {
	book, err := excelize.OpenReader(reader, opts...)
	if err != nil {
		return
	}
	r = &XlsxFile{book}
	return
}

func (r *XlsxFile) ReadSheet(num int, skipRows int, ch chan<- []any, cvfns ...data.ConvertFunc) {
	defer close(ch)
	name := r.GetSheetName(num)
	rows, _ := r.Rows(name)
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
func (r *XlsxFile) GetValues(num int) (values [][]string, err error) {
	name := r.GetSheetName(num)
	if name == "" {
		err = errors.New("sheet 不存在")
		return
	}
	return r.GetRows(name)
}

func Date(d float64) (date string, err error) {
	var t time.Time
	t, err = excelize.ExcelDateToTime(d, false)
	date = t.Format("2006-01-02 15:04:05")
	return
}
