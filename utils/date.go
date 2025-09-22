package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Weekdays 中文星期的表示
var Weekdays [7]string = [7]string{
	"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}

// IsLeapYear 判断指定的年份是否为润年
func IsLeapYear(year int) bool {
	return year%400 == 0 || year%4 == 0 && year%100 != 0
}

// GetMonthDays 获取指定月份的天数
func GetMonthDays(year, month int) (days int) {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		days = 31
	case 4, 6, 9, 11:
		days = 30
	case 2:
		if IsLeapYear(year) {
			days = 29
		} else {
			days = 28
		}
	default:
		panic("Wrong month")
	}
	return
}

// Date Basic type
type Date struct {
	time.Time
}

// Today 返回当前日期
func Today() Date {
	return Date{time.Now()}
}

// NewDate 构造日期
func NewDate(s string) (date Date) {
	patterns := [](*regexp.Regexp){
		regexp.MustCompile(`(\d{4})(\d{2})(\d{2})`),
		regexp.MustCompile(`(\d{4})[-/](\d{1,2})[-/](\d{1,2})`),
	}
	for _, r := range patterns {
		if r.MatchString(s) {
			result := r.FindStringSubmatch(s)
			var a [4]int
			var err error
			for i, v := range result[1:] {
				a[i], err = strconv.Atoi(v)
				if err != nil {
					panic("日期格式错")
				}
			}
			d := time.Date(a[0], time.Month(a[1]), a[2], 0, 0, 0, 0, time.Local)
			date = Date{d}
			return
		}
	}
	panic("日期不正确")
}

// Format 格式化日期
// %%  %
// %Y  4位年份，如：2020
// %M  2位月份，如：02
// %D  2位日期，如：15
// %Q  1位季度，如：2
// %W  星期，如星期一、星期二、星期日...
func (d Date) Format(format string) (s string) {
	f := func(s string) (r string) {
		switch s {
		case "%%":
			r = "%"
		case "%Y":
			r = fmt.Sprintf("%04d", d.Year())
		case "%y":
			r = fmt.Sprintf("%02d", d.Year()%100)
		case "%M":
			r = fmt.Sprintf("%02d", d.Month())
		case "%m":
			r = fmt.Sprintf("%2d", d.Month())
		case "%D":
			r = fmt.Sprintf("%02d", d.Day())
		case "%d":
			r = fmt.Sprintf("%2d", d.Day())
		case "%q":
			r = fmt.Sprintf("%1d", d.Quater())
		case "%Q":
			r = fmt.Sprintf("%1d季度", d.Quater())
		case "%F":
			r = fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
		case "%f":
			r = fmt.Sprintf("%04d%02d%02d", d.Year(), d.Month(), d.Day())
		case "%W":
			r = Weekdays[d.Weekday()]
		case "%w":
			r = fmt.Sprintf("%1d", d.Weekday())
		default:
			panic("Wrong Format")
		}
		return
	}
	r := regexp.MustCompile("%.")
	s = r.ReplaceAllStringFunc(format, f)
	return
}

// Add 增加 years,months,days
func (d Date) Add(years, months, days int) (date Date) {
	months = (d.Year()+years)*12 + int(d.Month()) + months - 1
	year := months / 12
	month := months%12 + 1
	day := min(GetMonthDays(year, month), d.Day())
	dat := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	date = Date{dat.AddDate(0, 0, days)}
	return
}

// Quater 返回日期的季度
func (d Date) Quater() int {
	return (int(d.Month()) + 2) / 3
}

// String 返回字符串
func (d Date) String() string {
	return d.Format("%F")
}

// NextDay 下一日
func (d Date) NextDay() Date {
	return d.Add(0, 0, 1)
}

// PrevDay 上一日
func (d Date) PrevDay() Date {
	return d.Add(0, 0, -1)
}

// TimeIt 统计一个函数运行所耗费的时间
// 使用方法，在函数的开始运行： defer TimeIt(time.Now())
func TimeIt(begin time.Time) {
	t := time.Since(begin)
	fmt.Println("本程序耗时：", t)
}
