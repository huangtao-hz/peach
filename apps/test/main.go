package main

import (
	"fmt"
	"peach/data"
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
	data := data.NewData()
	reader, err := book.NewReader("Sheet1", "A,C", 0)
	if err == nil {
		go reader.Read(data)
		data.Println()
	} else {
		fmt.Println(err)
	}

}
