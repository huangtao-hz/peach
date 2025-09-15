package main

import (
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
)

const create_sql = `
create table if not exists test(
	name    text,
	age     int
)
`

func main() {
	defer utils.Recover()

	db, err := sqlite.Open("test")
	utils.CheckFatal(err)
	defer db.Close()

	db.ExecScript(create_sql)
	sqlite.InitLoadFile(db)

	file := "~/abc.xlsx"
	r, err := excel.NewExcelReader(file)
	utils.CheckFatal(err)
	ch := make(chan []any)
	go r.ReadSheet(1, 0, ch)
	//utils.ChPrintln(data.Data)
	loader := sqlite.LoadFile(file, "test", ch)
	loader.FieldCount = 2
	loader.Load(db)
	db.Printf("select * from test", "%-12s   %6,d\n", "姓名             年龄", true)
	db.PrintRow("select * from test limit 1", "姓名,年龄")
}
