package main

import (
	"fmt"

	"peach/excel"
	"peach/utils"
	"time"
)

func writeexcel() {
	book := excel.NewWriter()
	sheet := book.GetSheet(0)
	sheet.Rename("Hello")

	//styles :=
	sheet.SetColStyle(map[string]string{
		"A":   "Short",
		"B:C": "Number",
		"D":   "Date",
	})
	sheet.SetWidth(map[string]float64{"A": 15,
		"B": 12,
		"C": 10,
		"D": 15,
	})

	sheet.AddTitle("A1:E1", "我的标题")
	sheet.AddHeader("A", "姓名,年龄,成绩,日期")
	sheet.AddRow("A", "张三", 12, 95, time.Now())
	sheet.AddRow("A", "王五", 12, 100, time.Now())
	sheet.SetBorder("A2:D4")
	sheet.SetColVisible("D", false)

	book.SaveAs("~/abcd.xlsx")
	fmt.Println("Sucess")
}

type Config struct {
	Name  string `toml:"name"`
	Home  string `toml:"home"`
	Hello string `toml:"hello"`
}

var s = `
[[sheet]]
A_D_E = 15.2

`

func main() {
	defer utils.Recover()
	//utils.PrintStruct(utils.Split("a|b|b|c  d"))
	//path := utils.NewPath("~/abc")
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
