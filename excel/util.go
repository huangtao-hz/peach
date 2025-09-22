package excel

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"os"
	"peach/data"
	"peach/utils"
	"slices"
)

// Book 定义工作簿接口
type Book interface {
	GetSheetList() []string
	IterRows(sheetIdx int, skipRows int) iter.Seq[[]string]
}

// ExcelBook 定义工作簿的类
type ExcelBook struct {
	Book
}

type ExcelReader interface {
	Read(sheets any, skipRows int, ch chan<- []any, cvfns ...data.ConvertFunc)
}

// NewExcelBook 构造行数
func NewExcelBook(reader io.Reader, path string) (book *ExcelBook, err error) {
	var _book Book
	ext := utils.NewPath(path).Ext()
	if ext == ".xls" {
		var (
			r  io.ReadSeeker
			ok bool
		)
		if r, ok = reader.(io.ReadSeeker); !ok {
			b, _ := io.ReadAll(reader)
			r = bytes.NewReader(b)
		}
		_book, err = NewXlsBook(r)
		if err != nil {
			return
		}
		book = &ExcelBook{_book}
	} else if slices.Contains([]string{".xlsx", ".xlsxm"}, ext) {
		_book, err = NewXlsxBook(reader)
		if err != nil {
			return
		}
		book = &ExcelBook{_book}
	}
	return
}

// GetSheets 根据指定的参数获取列表
func (b *ExcelBook) GetSheets(sheets any) (result []int) {
	sheetList := b.GetSheetList()
	if sheets == nil {
		result = make([]int, len(sheetList))
		for i := range sheetList {
			result[i] = i
		}
	} else {
		switch sheets := sheets.(type) {
		case int:
			result = []int{sheets}
		case string:
			i := slices.Index(sheetList, sheets)
			if i >= 0 {
				result = []int{i}
			}
		case []int:
			result = sheets
		case []string:
			result = make([]int, 0)
			for _, v := range sheets {
				i := slices.Index(sheetList, v)
				fmt.Println(i, v)
				if i >= 0 {
					result = append(result, i)
				}
			}
		}
	}
	return
}

// Read 读取 ExcelBook 文件内容
func (b *ExcelBook) Read(sheets any, skipRows int, ch chan<- []any, cvfns ...data.ConvertFunc) {
	defer close(ch)
	var convert = data.Convert(cvfns...)
	for _, idx := range b.GetSheets(sheets) {
		for row := range b.IterRows(idx, skipRows) {
			if r, err := convert(row); r != nil && err == nil {
				ch <- r
			}
		}
	}
}

// ExcelFile 定义 Excel
type ExcelFile struct {
	fp io.ReadSeekCloser
	ExcelBook
}

// NewExcelFile 构造函数
func NewExcelFile(path string) (f *ExcelFile, err error) {
	var (
		fp   io.ReadSeekCloser
		book Book
	)
	fp, err = os.Open(utils.Expand(path))
	if err != nil {
		return
	}
	book, err = NewExcelBook(fp, path)
	if err != nil {
		return
	}
	f = &ExcelFile{fp, ExcelBook{book}}
	return
}

// Close 关闭 ExcelReader
func (f *ExcelFile) Close() error {
	return f.fp.Close()
}
