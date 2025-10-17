package excel

import (
	"fmt"
	"slices"

	"github.com/huangtao-hz/excelize"
)

var NumberFmtMap = map[string]int{
	"Number":   3,
	"Currency": 4,
	"Percent":  10,
	"Time":     21,
	"Date":     28,
}

// Style 自定义的样式
type Style struct {
	FontName  string  // 字体名称
	FontSize  float64 // 字号
	Bold      bool    // 加粗
	Italic    bool    // 倾斜
	Underline string  // 下划线，为：single,double
	HAlign    string  // 水平对齐，为：left,center,right
	VAlign    string  // 垂直对齐，为：top,center,bottom
	WrapText  bool    // 字体换行
	NumFmt    string  // 数字样式
}

// AsStyle 转换为 excelize 支持的样式
func (s *Style) AsStyle() (style *excelize.Style, err error) {
	if s.Underline != "" && !slices.Contains([]string{"single", "double"}, s.Underline) {
		err = fmt.Errorf("%s 不是有效的 Underline 值", s.Underline)
		return
	}
	if s.HAlign != "" && !slices.Contains([]string{"left", "right", "center", "justify", "centerContinuous", "fill", "distributed"}, s.HAlign) {
		err = fmt.Errorf("%s 不是有效的 HAlign 值", s.HAlign)
		return
	}
	if s.VAlign != "" && !slices.Contains([]string{"top", "bottom", "center", "justify", "distributed"}, s.VAlign) {
		err = fmt.Errorf("%s 不是有效的 VAlign 值", s.VAlign)
		return
	}
	font := excelize.Font{
		Family:    s.FontName,
		Size:      s.FontSize,
		Bold:      s.Bold,
		Italic:    s.Italic,
		Underline: s.Underline,
	}
	aligment := excelize.Alignment{
		Horizontal: s.HAlign,
		Vertical:   s.VAlign,
		WrapText:   s.WrapText,
	}
	style = &excelize.Style{Font: &font, Alignment: &aligment}
	if s.NumFmt != "" {
		if numfmt, ok := NumberFmtMap[s.NumFmt]; ok {
			style.NumFmt = numfmt
		} else {
			style.CustomNumFmt = &s.NumFmt
		}
	}
	return
}

var predifinedStyles = map[string]Style{
	"Title": Style{
		FontName: "黑体",
		FontSize: 16,
		HAlign:   "centerContinuous",
		VAlign:   "center",
	},
	"Header": Style{
		FontName: "黑体",
		HAlign:   "center",
		VAlign:   "center",
	},
	"Percent": Style{
		NumFmt: "Percent",
		HAlign: "right",
		VAlign: "center",
	},
	"Currency": Style{
		NumFmt: "Currency",
		HAlign: "right",
		VAlign: "center",
	},
	"Number": Style{
		NumFmt: "Number",
		HAlign: "right",
		VAlign: "center",
	},
	"Date": Style{
		NumFmt: "Date",
		HAlign: "right",
		VAlign: "center",
	},
	"Normal": Style{
		HAlign:   "left",
		VAlign:   "center",
		WrapText: true,
	},
	"Short": Style{
		HAlign: "center",
		VAlign: "center",
	},
	"Normal-NoWrap": Style{
		HAlign: "left",
		VAlign: "center",
	},
}
