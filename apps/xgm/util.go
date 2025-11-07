package main

import (
	"embed"
	"fmt"
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

var (
	//go:embed tables
	tablesFS embed.FS
	//go:embed query
	queryFS embed.FS
	//go:embed loader
	loaderFS embed.FS
	//go:embed template
	templateFS embed.FS
)

// Reporter 报表类型
type Reporter struct {
	Sheet       string             `toml:"sheet"`
	Title       string             `toml:"title"`
	Header      string             `toml:"header"`
	Widths      map[string]float64 `toml:"widths"`
	Formats     map[string]string  `toml:"formats"`
	StartColumn string             `toml:"start_col"`
	Query       string             `toml:"query"`
	Hidden      string             `toml:"hidden"`

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
func (r *Reporter) Export(db *sqlite.DB, book *excel.Writer, args ...any) (err error) {
	if r.Sheet == "" {
		r.Sheet = "Sheet1"
	}
	sheet := book.GetSheet(r.Sheet)
	if r.Widths != nil {
		sheet.SetWidth(r.Widths)
	}
	if r.Formats != nil {
		sheet.SetColStyle(r.Formats)
	}
	if r.StartColumn == "" {
		r.StartColumn = "A"
	}
	if r.Hidden != "" {
		for col := range strings.SplitSeq(r.Hidden, ",") {
			sheet.SetColVisible(col, false)
		}
	}
	var rows *sqlite.Rows
	if rows, err = db.Query(r.Query, args...); err == nil {
		ch := make(chan []any, BufferSize)
		go rows.FetchAll(ch)
		if err = sheet.AddTable(r.StartColumn, r.Header, ch, r.Table); err == nil {
			sheet.SkipRows(2)
		}
	}
	return
}

// EpxortReport 导出报表
func ExportReport(db *sqlite.DB, book *excel.Writer, path string, args ...any) error {
	return NewReporter(path).Export(db, book, args...)
}

// ExportAll 导出多张报表
func ExportAll(db *sqlite.DB, book *excel.Writer, paths string) (err error) {
	for path := range strings.SplitSeq(paths, ",") {
		path = fmt.Sprintf("%s.toml", path)
		if err = (ExportReport(db, book, path)); err != nil {
			return
		}
	}
	return
}

// ExportXlsx 导出数据到 excel 文件中
func ExportXlsx(db *sqlite.DB, path string, files string) (err error) {
	w := excel.NewWriter()
	if err = ExportAll(db, w, files); err == nil {
		err = w.SaveAs(path)
	}
	return
}
