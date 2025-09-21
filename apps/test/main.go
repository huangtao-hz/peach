package main

import (
	"peach/excel"
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
	book, err := excel.NewExcelFile("~/Downloads/abc.xlsx")
	utils.CheckFatal(err)
	defer book.Close()
	ch := make(chan []any, 100)
	go book.Read(0, 0, ch, excel.UseCols("a:c"))
	utils.ChPrintln(ch)
}
