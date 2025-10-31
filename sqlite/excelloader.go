package sqlite

import (
	"errors"
	"io/fs"
	"peach/data"
	"peach/excel"
	"peach/utils"

	"github.com/BurntSushi/toml"
)

// ExcelLoader Excel 导入文件
type ExcelLoader struct {
	Loader
	Sheets   any    `tmol:"sheets"`
	UseCols  string `toml:"use_cols"`
	SkipRows int    `toml:"skip_rows"`
}

// NewExcelLoader 构造函数
func (db *DB) NewExcelLoader(fsys fs.FS, path string, book *excel.ExcelBook, fileinfo fs.FileInfo, cvfns ...data.ConvertFunc) (loader *ExcelLoader, err error) {
	loader = &ExcelLoader{Loader: Loader{db: db, Method: "insert", fileinfo: fileinfo, Clear: true, Check: true}}
	if _, err = toml.DecodeFS(fsys, path, loader); err != nil {
		return
	}
	if loader.Tablename == "" {
		err = errors.New("tablename can't be empty.")
		return
	}
	loader.reader, err = book.NewReader(loader.Sheets, loader.UseCols, loader.SkipRows, cvfns...)
	return
}

// LoadExcel 导入excel文件
func (db *DB) LoadExcel(fsys fs.FS, path string, book *excel.ExcelBook, fileinfo fs.FileInfo, cvfns ...data.ConvertFunc) error {
	if loader, err := db.NewExcelLoader(fsys, path, book, fileinfo, cvfns...); err == nil {
		return loader.Load()
	}
	return nil
}

// LoadExcelFile  导入单个 Excel 文件
func (db *DB) LoadExcelFile(fsys fs.FS, path string, file utils.File, cvfns ...data.ConvertFunc) error {
	f, err := excel.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return db.LoadExcel(fsys, path, &f.ExcelBook, file.FileInfo(), cvfns...)
}
