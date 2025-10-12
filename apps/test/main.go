package main

import (
	"fmt"
	"peach/archive"
	"peach/excel"
	"peach/utils"
	"time"

	"github.com/huangtao-hz/excelize"
)

func proc(file archive.File) {
	if !file.FileInfo().IsDir() {
		fmt.Println("File:", file.FileInfo().Name())
		if book, err := excel.Open(file); err == nil {
			fmt.Println(book.GetSheetList())
		} else {
			fmt.Println("Error:", err)
		}
	}
}
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
	sheet.AddHeader("A", 2, "姓名,年龄,成绩,日期")
	sheet.AddRow("A", 3, "张三", 12, 95, time.Now())
	sheet.AddRow("A", 4, "王五", 12, 100, time.Now())
	sheet.SetBorder("A2:D4")
	sheet.SetColVisible("D", false)

	book.SaveAs("~/abcd.xlsx")
	fmt.Println("Sucess")
}
func archivetest() {
	file := utils.NewPath("~/Downloads").Find("20250905.tgz")
	if file != nil {
		fmt.Println(file.String())
		//fmt.Println(file.FileInfo().Name())
		archive.ExtractTar(file.String(), proc)
	}
}

func main() {
	defer utils.Recover()
	//utils.PrintStruct(utils.Split("a|b|b|c  d"))
	book := excel.NewWriter()
	sheet := book.GetSheet(0)
	sheet.AddHeader("A", 1, "姓名,年龄,学校")
	sheet.AddRow("A", 2, "张三", 18, "黔江司机称")
	sheet.AddRow("A", 3, "王五", 20, "河南附中")
	book.AddTable("Sheet1", &excelize.Table{
		Range:     "A1:C3",
		StyleName: "TableStyleMedium6",
	})

	book.SaveAs("~/abc.xlsx")
}
