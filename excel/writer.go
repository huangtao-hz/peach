package excel

import (
	"fmt"
	"peach/utils"
	"strings"

	"github.com/xuri/excelize/v2"
)

// 把坐标转换成单位元格
func Cell(col any, row int) (cell string) {
	var err error
	switch column := col.(type) {
	case int:
		cell, err = excelize.CoordinatesToCellName(column, row)
		utils.CheckFatal(err)
	case string:
		cell, err = excelize.JoinCellName(column, row)
		utils.CheckFatal(err)
	}
	return
}

// Excel 文件
type File struct {
	*excelize.File
}

// 新建 Excel 文件
func NewFile() *File {
	return &File{excelize.NewFile()}
}

// 获取工作表
func (f *File) GetSheet(index any) *WorkSheet {
	var name string
	switch idx := index.(type) {
	case int:
		name = f.GetSheetName(idx)
	case string:
		name = idx
		i, _ := f.GetSheetIndex(name)
		if i == -1 {
			f.NewSheet(name)
		}
	}
	return &WorkSheet{f, name}
}

// 保存文件
func (f *File) SaveAs(path string) {
	f.File.SaveAs(utils.Expand(path))
}

// 工作表
type WorkSheet struct {
	file *File
	name string
}

// 设置表格的宽度
func (s *WorkSheet) SetWidth(widthes map[string]float64) {
	for col, width := range widthes {
		aa := strings.Split(col, ":")
		if len(aa) == 1 {
			aa = append(aa, aa[0])
		}
		if len(aa) != 2 {
			panic("单元格格式错")
		}
		s.file.SetColWidth(s.name, aa[0], aa[1], width)
	}
}

// 修改工作表名称
func (s *WorkSheet) Rename(newName string) {
	s.file.SetSheetName(s.name, newName)
	s.name = newName
}

// TableFormat 单位元样式
const TableFormat = `{"table_style":"TableStyleMedium6", "show_first_column":false,"show_last_column":false,"show_row_stripes":true,"show_column_stripes":false}`

// 写入表格
func (s *WorkSheet) WriteTable(axis string, header any, ch <-chan []any) {
	var count int
	var table *excelize.Table
	col, row, err := excelize.CellNameToCoordinates(axis) // 读取初始
	utils.CheckFatal(err)
	writer, err := s.file.NewStreamWriter(s.name)
	utils.CheckFatal(err)
	if header != nil {
		var headers []string
		switch head := header.(type) {
		case []string:
			headers = head
		case string:
			headers = strings.Split(head, ",")
		default:
			panic("标题必须为 string 或 []string")
		}
		count = len(headers)
		writer.SetRow(Cell(col, row), utils.Slice(headers))
		row++
	} else {
		panic("未设置标题")
	}
	for rowdata := range ch {
		writer.SetRow(Cell(col, row), rowdata)
		row++
	}
	end, _ := excelize.CoordinatesToCellName(col+count-1, row-1)
	table = &excelize.Table{}
	table.Range = fmt.Sprintf("%s:%s", axis, end)
	table.StyleName = "TableStyleMedium6"
	writer.AddTable(table)
	writer.Flush()
}
