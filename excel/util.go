package excel

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	dateFormat1 *regexp.Regexp = regexp.MustCompile(`^\d{2}-\d{2}-\d{2}$`)
	dateFormat2 *regexp.Regexp = regexp.MustCompile(`^\d+$`)
	dateFormat3 *regexp.Regexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
)

// FormatDate 格式化 Excel 文件中的日期
// 将日期统格式化成 YYYY-MM-DD 格式，
// 无法格式化的，沿用原来的值
func FormatDate(s string, layout string) string {
	if d, err := AToDate(s, layout); err == nil {
		return d
	}
	return s
}

// AToDate 格式化日期
func AToDate(s string, layout string) (d string, err error) {
	var (
		t time.Time
		f float64
	)
	if dateFormat1.MatchString(s) {
		t, err = time.Parse("01-02-06", s)
	} else if dateFormat2.MatchString(s) {
		if f, err = strconv.ParseFloat(s, 64); err == nil {
			t, err = excelize.ExcelDateToTime(f, false)
		}
	} else if dateFormat3.MatchString(s) {
		t, err = time.Parse("2006-01-02", s)
	} else {
		err = fmt.Errorf("无法将 %v 转换成日期", s)
	}
	if err == nil {
		d = t.Format(layout)
	}
	return
}
