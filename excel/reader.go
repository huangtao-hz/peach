package excel

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	dateFormat1 *regexp.Regexp = regexp.MustCompile(`^\d{2}-\d{2}-\d{2}$`)
	dateFormat2 *regexp.Regexp = regexp.MustCompile(`^\d+$`)
)

// FormatDate 格式化 Excel 文件中的日期
// 将日期统格式化成 YYYY-MM-DD 格式，
// 无法格式化的，沿用原来的值
func FormatDate(s string) string {
	if dateFormat1.MatchString(s) {
		h := strings.Split(s, "-")
		return fmt.Sprintf("20%s-%s-%s", h[2], h[0], h[1])
	}
	if dateFormat2.MatchString(s) {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			d, _ := Date(f)
			return d[:10]
		}
	}
	return s
}
