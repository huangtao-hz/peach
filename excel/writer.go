package excel

import (
	_ "embed"
	"fmt"
	"peach/utils"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/huangtao-hz/excelize"
)

const (
	InnerBorderStyle = 1
	OuterBorderStyle = 2
	BorderColor      = "5F5C5D"
)

// Writer Excel 文件
type Writer struct {
	*excelize.File
	styles map[string]int
	sheets map[string]*WorkSheet
}

//go:embed styles.toml
var prestyles string

// NewFile 新建 Excel 文件
func NewWriter() (file *Writer) {
	f := excelize.NewFile()
	styles := make(map[string]int)
	file = &Writer{File: f, styles: styles}
	utils.CheckErr(file.AddPreStyles(prestyles))
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
			w.styles[name] = id
		}
	}
	return
}

// GetSheet 获取工作表
func (w *Writer) GetSheet(name string) *WorkSheet {
	if w.sheets == nil {
		w.SetSheetName("Sheet1", name)
		w.sheets = make(map[string]*WorkSheet)
	} else if st, ok := w.sheets[name]; ok {
		return st
	} else {
		w.NewSheet(name)
	}
	st := &WorkSheet{writer: w, name: name, Row: 1}
	w.sheets[name] = st
	return st
}

// SaveAs 保存文件
func (w *Writer) SaveAs(path string) error {
	return w.File.SaveAs(utils.Expand(path))
}

// WorkSheet 工作表
type WorkSheet struct {
	writer *Writer
	name   string
	Row    int
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
			if styleID, ok := s.writer.styles[styleName]; !ok {
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

// AddTableToml 使用 Toml 来描述表格
func (s *WorkSheet) AddTableToml(cell string, table string, ch <-chan []any) (err error) {
	table_ := &excelize.Table{}
	toml.Decode(table, table_)
	//utils.PrintStruct(table_)
	return s.AddTable(cell, "", ch, table_)
}

// AddTable 写入表格
func (s *WorkSheet) AddTable(axis string, header string, ch <-chan []any, opt ...*excelize.Table) (err error) {
	var (
		row, col  int
		table     *excelize.Table
		total_row bool
	)
	if opt == nil || opt[0] == nil {
		table = &excelize.Table{}
	} else {
		table = opt[0]
	}

	if table.ShowHeaderRow == nil {
		showHeaderRow := true
		table.ShowHeaderRow = &showHeaderRow
	}
	col, row, err = excelize.CellNameToCoordinates(axis) // 读取初始
	if err != nil {
		if col, err = excelize.ColumnNameToNumber(axis); err == nil {
			row = s.Row
			axis, _ = excelize.JoinCellName(axis, row)
		} else {
			return
		}
	} else {
		s.Row = row
	}

	colname, _ := excelize.ColumnNumberToName(col)
	if *table.ShowHeaderRow {
		s.Row++
	}
	for rowdata := range ch {
		s.AddRow(colname, rowdata...)
	}
	count := slices.Max([]int{len(utils.Split(header)), len(table.Columns)})
	end, _ := excelize.CoordinatesToCellName(col+count-1, s.Row-1)

	table.Range = fmt.Sprintf("%s:%s", axis, end)
	if table.StyleName == "" {
		table.StyleName = "TableStyleMedium6"
	}
	var proc_style = func(style *string) {
		if *style != "" {
			if styleId, ok := s.writer.styles[*style]; ok {
				*style = fmt.Sprint(styleId)
			} else {
				*style = ""
			}
		}
	}

	if table.Columns != nil {
		for i := range len(table.Columns) {
			column := &table.Columns[i]
			if column.HeaderRowCellStyle == "" {
				column.HeaderRowCellStyle = "Header"
			}
			proc_style(&column.DataCellStyle)
			proc_style(&column.TotalsRowCellStyle)
			proc_style(&column.HeaderRowCellStyle)
			if column.TotalsRowCellStyle == "" && column.DataCellStyle != "" {
				column.TotalsRowCellStyle = column.DataCellStyle
			}
			if column.TotalsRowFunction != "" || column.TotalsRowLabel != "" {
				total_row = true
			}
		}
	} else if header != "" {
		columns := make([]excelize.TableColumn, 0)
		for _, h := range utils.Split(header) {
			column := excelize.TableColumn{Name: h, HeaderRowCellStyle: "Header"}
			proc_style(&column.HeaderRowCellStyle)
			columns = append(columns, column)
		}
		table.Columns = columns
	}

	s.writer.AddTable(s.name, table)
	if total_row {
		s.Row++
	}
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
		if styleName != "" {
			return s.SetCellStyle(rng, styleName)
		}
	} else {
		return err
	}
	return nil
}

// SetCellStyle 设置单元格样式
func (s *WorkSheet) SetCellStyle(rng string, styleName string) (err error) {
	var cell1, cell2 string
	if cell1, cell2, err = Range2Cells(rng); err == nil {
		if styleID, ok := s.writer.styles[styleName]; ok {
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
		col1, row1, col2, row2 int
		style                  *excelize.Style
		id                     int
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
			left_style, right_style, top_style, bottom_style := InnerBorderStyle, InnerBorderStyle, InnerBorderStyle, InnerBorderStyle
			if c == col1 {
				left_style = OuterBorderStyle
			}
			if c == col2 {
				right_style = OuterBorderStyle
			}
			if r == row1 {
				top_style = OuterBorderStyle
			}
			if r == row2 {
				bottom_style = OuterBorderStyle
			}
			style.Border = []excelize.Border{
				excelize.Border{Type: "left", Style: left_style, Color: BorderColor},
				excelize.Border{Type: "right", Style: right_style, Color: BorderColor},
				excelize.Border{Type: "top", Style: top_style, Color: BorderColor},
				excelize.Border{Type: "bottom", Style: bottom_style, Color: BorderColor},
			}
			if id, err = s.writer.NewStyle(style); err != nil {
				return
			}
			s.writer.SetCellStyle(s.name, cell, cell, id)
		}
	}
	return
}

// SkipRows 跳过空行
func (s *WorkSheet) SkipRows(rows int) (row int) {
	s.Row += rows
	return s.Row
}

// get_cell 获取当前坐标
func (s *WorkSheet) get_cell(col string) (cellname string, err error) {
	return excelize.JoinCellName(col, s.Row)
}

// AddTitle 添加标题
func (s *WorkSheet) AddTitle(rng string, title string) (err error) {
	cells := strings.Split(rng, ":")
	if len(cells) == 0 {
		return fmt.Errorf("%s is not a valid range format", rng)
	} else if len(cells) == 1 {
		cells = append(cells, cells[0])
	}
	var cell1, cell2 string
	if cell1, err = s.get_cell(cells[0]); err == nil {
		if err = s.SetCellValue(cell1, title); err == nil {
			if cell2, err = s.get_cell(cells[1]); err == nil {
				if err = s.SetCellStyle(strings.Join([]string{cell1, cell2}, ":"), "Title"); err == nil {
					s.SetRowHeight(s.Row, 30)
					s.Row++
				}
			}
		}
	}
	return
}

// AddHeader 添加标题,header 是用 “,” 分隔的字符串
func (s *WorkSheet) AddHeader(col string, header string) (err error) {
	var col_, row int
	row = s.Row
	if col_, err = excelize.ColumnNameToNumber(col); err == nil {
		for i, h := range strings.Split(header, ",") {
			if err = s.SetCell(col_+i, row, h, "Header"); err != nil {
				return err
			}
		}
	}
	s.Row++
	return
}

// AddHeader 添加一行数据
func (s *WorkSheet) AddRow(col string, values ...any) (err error) {
	var col_, row int
	row = s.Row
	if col_, err = excelize.ColumnNameToNumber(col); err == nil {
		for i, value := range values {
			if err = s.SetCell(col_+i, row, value, ""); err != nil {
				return err
			}
		}
	}
	s.Row++
	return
}
