package main

import (
	"fmt"

	"peach/excel"
	"peach/utils"
)

type Config struct {
	Name  string `toml:"name"`
	Home  string `toml:"home"`
	Hello string `toml:"hello"`
}

var s = `

[widths]
"D,E" = 15.2
"A:B" = 13.2
`

type Abc struct {
	Widths map[string]float64 `toml:"widths"`
}

func test_xls() {
	file := excel.NewWriter()
	defer file.SaveAs("~/Documents/abc.xlsx")
	sheet := file.GetSheet("test")
	sheet.AddTitle("A:F", "这是一个大标题")
	sheet.AddHeader("A", "我们,他们,历史,政治,音乐,天空")
	sheet.AddRow("A", "a", 12, 23, 34, 46)
	sheet.AddRow("A", "b", 24, 32, 23, 146)
	sheet.SetBorder(fmt.Sprintf("A2:F%d", sheet.Row-1))
	fmt.Println("生成文件成功！")
}

func main() {
	defer utils.Recover()
	//utils.PrintStruct(utils.Split("a|b|b|c  d"))
	//path := utils.NewPath("~/abc")

}
