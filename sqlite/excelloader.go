package sqlite

import (
	"io"
	"io/fs"
	"peach/data"
	"peach/excel"
	"peach/utils"

	"github.com/BurntSushi/toml"
)

// ExcelLoader Excel 导入文件
type excelLoader struct {
	Loader
	Sheets   any    `tmol:"sheets"`
	UseCols  string `toml:"use_cols"`
	SkipRows int    `toml:"skip_rows"`
}

// SingleExcelLoader
type SingleExcelLoader struct {
	file io.ReadCloser
	excelLoader
}

// NewExcelLoader 构造函数
func (db *DB) NewExcelLoader(fsys fs.FS, path string, excelFile utils.File, cvfns ...data.ConvertFunc) (loader *SingleExcelLoader, err error) {
	loader = &SingleExcelLoader{}
	if _, err = toml.DecodeFS(fsys, path, loader); err != nil {
		return
	}
	var (
		book *excel.ExcelBook
	)
	if loader.file, err = excelFile.Open(); err != nil {
		return
	}
	if book, err = excel.NewExcelBook(loader.file, excelFile.FileInfo().Name()); err != nil {
		return
	}
	loader.reader, err = book.NewReader(loader.Sheets, loader.UseCols, loader.SkipRows, cvfns...)
	return
}

// Close 关闭 Excel 文件
func (l *SingleExcelLoader) Close() error {
	return l.file.Close()
}

// ExcelLoader
type ExcelLoader struct {
	file   io.ReadCloser
	loders []*excelLoader
}
