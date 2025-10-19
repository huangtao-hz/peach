package excel

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"peach/data"
	"peach/utils"
	"slices"
)

// Book 定义工作簿接口
type book interface {
	GetSheetList() []string
	IterRows(sheetIdx int, skipRows int) iter.Seq[[]string]
}

// ExcelBook 定义工作簿的类
type ExcelBook struct {
	book
}

// NewExcelBook 构造行数
func NewExcelBook(reader io.Reader, path string) (ebook *ExcelBook, err error) {
	var _book book
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
		_book, err = newXlsBook(r)
		if err != nil {
			return
		}
		ebook = &ExcelBook{_book}
	} else if slices.Contains([]string{".xlsx", ".xlsxm"}, ext) {
		_book, err = newXlsxBook(reader)
		if err != nil {
			return
		}
		ebook = &ExcelBook{_book}
	} else {
		err = fmt.Errorf("%s 不是 excel 文件", path)
	}
	return
}

// GetSheets 根据指定的参数获取列表
func (b *ExcelBook) GetSheets(sheets any) (result []int, err error) {
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
				if i >= 0 {
					result = append(result, i)
				}
			}
		}
	}
	return
}

// NewReader 新建 ExcelReader
func (b *ExcelBook) NewReader(sheets any, useCols string, skipRows int, cvfns ...data.ConvertFunc) (r *Reader, err error) {
	var sheetlist []int
	if useCols != "" {
		cvfns = slices.Insert(cvfns, 0, UseCols(useCols))
	}
	if sheetlist, err = b.GetSheets(sheets); err == nil {
		r = &Reader{b, sheetlist, skipRows, cvfns}
	}
	return
}

// ExcelFile 定义 Excel
type ExcelFile struct {
	fp io.ReadCloser
	ExcelBook
}

// NewExcelFile 构造函数
func NewExcelFile(path string) (f *ExcelFile, err error) {
	return Open(utils.NewPath(path))
}

// Open 打开 Excel 文件
func Open(path utils.File) (f *ExcelFile, err error) {
	var (
		fp   io.ReadCloser
		book book
	)
	if fp, err = path.Open(); err != nil {
		return
	}
	if book, err = NewExcelBook(fp, path.FileInfo().Name()); err != nil {
		return
	}
	f = &ExcelFile{fp, ExcelBook{book}}
	return
}

// Close 关闭 ExcelReader
func (f *ExcelFile) Close() error {
	return f.fp.Close()
}

// Reader 定义读取 Excel 的机构体
type Reader struct {
	book     *ExcelBook
	sheets   []int
	SkipRows int
	cvfns    []data.ConvertFunc
}

// Read 读取 Excel 数据
func (r *Reader) Read(d *data.Data) {
	defer close(d.Data)
	var convert = data.Convert(r.cvfns...)
	for _, idx := range r.sheets {
		for row := range r.book.IterRows(idx, r.SkipRows) {
			if r, err := convert(row); err != nil {
				d.Cancel(err)
			} else if r != nil {
				select {
				case d.Data <- r:
				case <-d.Done():
					return
				}
			}
		}
	}
}
