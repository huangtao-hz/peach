package main

import (
	"embed"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	BufferSize = 1024
)

//go:embed tables
var tablesFS embed.FS

// Reporter 报表类型
type Reporter struct {
	Title       string `toml:"title"`
	Header      string `toml:"header"`
	StartColumn string `toml:"start_col"`
	Query       string `toml:"query"`
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
func (r *Reporter) Export(db *sqlite.DB, sheet *excel.WorkSheet, args ...any) error {
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
func ExportReport(db *sqlite.DB, sheet *excel.WorkSheet, path string, args ...any) error {
	return NewReporter(path).Export(db, sheet)
}
