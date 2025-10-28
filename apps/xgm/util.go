package main

import (
	"embed"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/huangtao-hz/excelize"
)

const (
	BufferSize = 1024
)

//go:embed tables
var tablesFS embed.FS

// Reporter 报表类型
type Reporter struct {
	Sheet       string             `toml:"sheet"`
	Title       string             `toml:"title"`
	Header      string             `toml:"header"`
	Widths      map[string]float64 `toml:"widths"`
	Formats     map[string]string  `toml:"formats"`
	StartColumn string             `toml:"start_col"`
	Query       string             `toml:"query"`

	*excelize.Table
}

// NewReporter 构造函数
func NewReporter(path string) *Reporter {
	rep := Reporter{}
	if !strings.HasPrefix(path, "tables/") {
		path = strings.Join([]string{"tables", path}, "/")
	}
	_, err := toml.DecodeFS(tablesFS, path, &rep)
	utils.CheckFatal(err)
	return &rep
}

// Export 导出报表
func (r *Reporter) Export(db *sqlite.DB, book *excel.Writer, args ...any) error {
	sheet := book.GetSheet(r.Sheet)
	if r.Widths != nil {
		sheet.SetWidth(r.Widths)
	}
	if r.Formats != nil {
		sheet.SetColStyle(r.Formats)
	}
	if rows, err := db.Query(r.Query, args...); err == nil {
		ch := make(chan []any, BufferSize)
		go rows.FetchAll(ch)
		if r.Header != "" {
			return sheet.AddTable(r.StartColumn, r.Header, ch)
		}
	} else {
		return err
	}
	return nil
}

// EpxortReport 导出报表
func ExportReport(db *sqlite.DB, book *excel.Writer, path string, args ...any) error {
	return NewReporter(path).Export(db, book)
}
