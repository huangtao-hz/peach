package sqlite

import (
	"io/fs"
	"peach/data"
	"peach/excel"

	"github.com/BurntSushi/toml"
)

// ExcelLoader Excel 导入文件
type ExcelLoader struct {
	Loader
	Sheets   any    `tmol:"sheets"`
	UseCols  string `toml:"use_cols"`
	SkipRows int    `toml:"skip_rows"`
}

type ExcelLoaders struct {
	Loaders []*ExcelLoader `toml:"loader"`
}

// NewExcelLoader 构造函数
func (db *DB) NewExcelLoader(fsys fs.FS, path string, book *excel.ExcelBook, cvfns ...data.ConvertFunc) (loader *ExcelLoader, err error) {
	loader = &ExcelLoader{}
	if _, err = toml.DecodeFS(fsys, path, loader); err != nil {
		return
	}
	loader.reader, err = book.NewReader(loader.Sheets, loader.UseCols, loader.SkipRows, cvfns...)
	return
}
