package excel

import (
	"fmt"
	"peach/utils"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	InnerBorderStyle = 1
	OuterBorderStyle = 2
	BorderColor      = "5F5C5D"
)

// Writer Excel 文件
type Writer struct {
	*excelize.File
	Styles map[string]int
}

// NewFile 新建 Excel 文件
func NewWriter() (file *Writer) {
	f := excelize.NewFile()
	styles := make(map[string]int)
	file = &Writer{File: f, Styles: styles}
	for name, style := range predifinedStyles {
		if err := file.AddStyle(name, style); err != nil {
			fmt.Println(err)
		}
	}
	return
}

// AddStyle 新增样式
func (w *Writer) AddStyle(name string, style Style) (err error) {
	var (
		id int
		s  *excelize.Style
	)
	if s, err = style.AsStyle(); err == nil {
		if id, err = w.NewStyle(s); err == nil {
			w.Styles[name] = id
		}
	}
	return
}

// GetSheet 获取工作表
func (w *Writer) GetSheet(index any) *WorkSheet {
	var name string
	switch idx := index.(type) {
	case int:
		name = w.GetSheetName(idx)
	case string:
		name = idx
		i, _ := w.GetSheetIndex(name)
		if i == -1 {
			w.NewSheet(name)
		}
	}
	return &WorkSheet{w, name}
}

// SaveAs 保存文件
func (w *Writer) SaveAs(path string) {
	w.File.SaveAs(utils.Expand(path))
}

// WorkSheet 工作表
type WorkSheet struct {
	writer *Writer
	name   string
}

// SetWidth 设置表格的宽度
func (s *WorkSheet) SetWidth(widthes map[string]float64) (err error) {
	for cols, width := range widthes {
		for col := range strings.SplitSeq(cols, ",") {
			aa := strings.Split(col, ":")
			if len(aa) == 1 {
				aa = append(aa, aa[0])
			} else if len(aa) != 2 {
				return fmt.Errorf("%s 不是有效的单元格格式", col)
			}
			s.writer.SetColWidth(s.name, aa[0], aa[1], width)
		}
	}
	return
}

// SetColStyle 设置列样式
func (s *WorkSheet) SetColStyle(styles map[string]string) error {
	for columns, styleName := range styles {
		for column := range strings.SplitSeq(columns, ",") {
			if styleID, ok := s.writer.Styles[styleName]; !ok {
				return fmt.Errorf("样式 %s 不存在", styleName)
			} else {
				s.writer.SetColStyle(s.name, column, styleID)
			}
		}
	}
	return nil
}

// SetColHidden
func (s *WorkSheet) SetColVisible(col string, visible bool) error {
	return s.writer.SetColVisible(s.name, col, visible)
}

// Rename 修改工作表名称
func (s *WorkSheet) Rename(newName string) {
	s.writer.SetSheetName(s.name, newName)
	s.name = newName
}

// tableFormat 单位元样式
const tableFormat = `{"table_style":"TableStyleMedium6", "show_first_column":false,"show_last_column":false,"show_row_stripes":true,"show_column_stripes":false}`

// WriteTable 写入表格
func (s *WorkSheet) WriteTable(axis string, header string, ch <-chan []any) (err error) {
	var row, col int
	col, row, err = excelize.CellNameToCoordinates(axis) // 读取初始
	if err != nil {
		return
	}
	colname, _ := excelize.ColumnNumberToName(col)
	s.AddHeader(colname, row, header)
	row++
	for rowdata := range ch {
		s.AddRow(colname, row, rowdata...)
		row++
	}
	count := len(strings.Split(header, ","))
	end, _ := excelize.CoordinatesToCellName(col+count-1, row-1)
	table := &excelize.Table{
		Range:     fmt.Sprintf("%s:%s", axis, end),
		StyleName: "TableStyleMedium6",
	}
	s.writer.AddTable(s.name, table)
	return
}

// SetCellValue 设置单元格值
func (s *WorkSheet) SetCellValue(cell string, value any) error {
	return s.writer.SetCellValue(s.name, cell, value)
}

// MergeCell 合并单元格
func (s *WorkSheet) MergeCell(topleft string, bottomright string) error {
	return s.writer.MergeCell(s.name, topleft, bottomright)
}

// SetRowHeight 设置行高
func (s *WorkSheet) SetRowHeight(row int, height float64) error {
	return s.writer.SetRowHeight(s.name, row, height)
}

// GetCellStyle 获取单元格样式
func (s *WorkSheet) GetCellStyle(cell string) (style *excelize.Style, err error) {
	var id int
	if id, err = s.writer.GetCellStyle(s.name, cell); err == nil {
		return s.writer.GetStyle(id)
	}
	return
}

// SetMergeCell 设置合并单元格
func (s *WorkSheet) SetMergeCell(rng string, value any, styleName string) error {
	if cell1, cell2, err := Range2Cells(rng); err == nil {
		s.MergeCell(cell1, cell2)
		s.SetCellValue(cell1, value)
		return s.SetCellStyle(rng, styleName)
	} else {
		return err
	}
}

// SetCellStyle 设置单元格样式
func (s *WorkSheet) SetCellStyle(rng string, styleName string) (err error) {
	var cell1, cell2 string
	if cell1, cell2, err = Range2Cells(rng); err == nil {
		if styleID, ok := s.writer.Styles[styleName]; ok {
			s.writer.SetCellStyle(s.name, cell1, cell2, styleID)
		} else {
			err = fmt.Errorf("样式 %s 不存在", styleName)
		}
	}
	return
}

// SetCell 设置单元格
func (s *WorkSheet) SetCell(col any, row int, value any, styleName string) (err error) {
	var cell string
	if cell, err = Cell(col, row); err == nil {
		if err = s.SetCellValue(cell, value); err == nil {
			if styleName != "" {
				return s.SetCellStyle(cell, styleName)
			}
		}
	}
	return
}

// SetBorder 设置边框
func (s *WorkSheet) SetBorder(rng string) (err error) {
	var (
		col1, row1, col2, row2   int
		style                    *excelize.Style
		left, right, top, bottom excelize.Border
		id                       int
	)
	if col1, row1, col2, row2, err = RangeToCoordinates(rng); err != nil {
		return
	}
	for c := col1; c <= col2; c++ {
		for r := row1; r <= row2; r++ {
			cell, _ := Cell(c, r)
			if style, err = s.GetCellStyle(cell); err != nil {
				return
			}
			if c == col1 {
				left = excelize.Border{Type: "left", Style: OuterBorderStyle, Color: BorderColor}
			} else {
				left = excelize.Border{Type: "left", Style: InnerBorderStyle, Color: BorderColor}
			}
			if c == col2 {
				right = excelize.Border{Type: "right", Style: OuterBorderStyle, Color: BorderColor}
			} else {
				right = excelize.Border{Type: "right", Style: InnerBorderStyle, Color: BorderColor}
			}
			if r == row1 {
				top = excelize.Border{Type: "top", Style: OuterBorderStyle, Color: BorderColor}
			} else {
				top = excelize.Border{Type: "top", Style: InnerBorderStyle, Color: BorderColor}
			}
			if r == row2 {
				bottom = excelize.Border{Type: "bottom", Style: OuterBorderStyle, Color: BorderColor}
			} else {
				bottom = excelize.Border{Type: "bottom", Style: InnerBorderStyle, Color: BorderColor}
			}
			style.Border = []excelize.Border{left, top, right, bottom}
			if id, err = s.writer.NewStyle(style); err != nil {
				return
			}
			s.writer.SetCellStyle(s.name, cell, cell, id)
		}
	}
	return
}

// AddTitle 添加标题
func (s *WorkSheet) AddTitle(rng string, title string) (err error) {
	var cell1 string
	if cell1, _, err = Range2Cells(rng); err != nil {
		return
	}
	if err = s.SetCellValue(cell1, title); err != nil {
		return
	}
	return s.SetCellStyle(rng, "Title")
}

// AddHeader 添加标题,header 是用 “,” 分隔的字符串
func (s *WorkSheet) AddHeader(col string, row int, header string) (err error) {
	var col_ int
	if col_, err = excelize.ColumnNameToNumber(col); err == nil {
		for i, h := range strings.Split(header, ",") {
			if err = s.SetCell(col_+i, row, h, "Header"); err != nil {
				return err
			}
		}
	}
	return
}

// AddHeader 添加一行数据
func (s *WorkSheet) AddRow(col string, row int, values ...any) (err error) {
	var col_ int
	if col_, err = excelize.ColumnNameToNumber(col); err == nil {
		for i, value := range values {
			if err = s.SetCell(col_+i, row, value, ""); err != nil {
				return err
			}
		}
	}
	return
}
