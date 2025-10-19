package main

type ExcelLoader struct {
	Table      string `toml:"table"`
	Method     string `toml:"method"`
	Fields     string `toml:"fields"`
	FieldCount int    `toml:"field_count"`
	LoadSQL    string `toml:"load_sql"`
	Check      string `toml:"check"`
	Clear      string `toml:"clear"`
	ClearSQL   string `toml:"clear_sql"`
	Sheets     any    `toml:"Sheets"`
	UseCols    string `toml:"use_cols"`
	SkipRows   int    `toml:"skip_rows"`
}
