package main

import (
	"embed"
	"peach/excel"
	"peach/sqlite"
	"strings"

	"github.com/BurntSushi/toml"
)

//go:embed tables
var tablesFS embed.FS

type Reporter struct {
	Title       string `toml:"title"`
	Header      string `toml:"header"`
	StartColumn string `toml:"start_col"`
	Query       string `toml:"query"`
}

func NewReporter(path string) (*Reporter, error) {
	rep := Reporter{}
	if !strings.HasPrefix(path, "tables/") {
		path = strings.Join([]string{"tables", path}, "/")
	}
	_, err := toml.DecodeFS(tablesFS, path, &rep)
	return &rep, err
}

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

func test_export(db *sqlite.DB) error {
	file := excel.NewWriter()
	defer file.SaveAs("~/Documents/abc.xlsx")
	sheet := file.GetSheet("test")
	if r, err := NewReporter("tcjyb.toml"); err == nil {
		r.Export(db, sheet)
	} else {
		return err
	}
	return nil
}
