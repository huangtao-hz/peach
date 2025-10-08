// Package 	: utils
// Writer	: Huangtao
// Create	: 2025-08-05

package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CheckFatal 检查系统是否有致命错误
// 如有则退出系统
func CheckFatal(err error) {
	if err != nil {
		panic(err)
	}
}

// Recover 捕获错误并提示，避免程序崩溃
// 在主程序开始时使用 defer 调用
func Recover() {
	if err := recover(); err != nil {
		fmt.Println("Fatal:", err)
	}
}

// Slice 将任意类型转换为 []any
func Slice[T any](s []T) (d []any) {
	var ok bool
	if d, ok = any(s).([]any); ok {
		return
	} else {
		d = make([]any, len(s))
		for i, v := range s {
			d[i] = v
		}
	}
	return
}

// formatInt 对整数进行格式化，增加千分节
func formatInt(k string) string {
	length := len(k)
	c, m := length/3, length%3
	b := make([]string, 0)
	if m > 0 {
		b = append(b, k[:m])
	}
	for i := range c {
		b = append(b, k[i*3+m:i*3+m+3])
	}
	return strings.Replace(strings.Join(b, ","), "-,", "-", 1)
}

// Sprintf 字符串格式化，解决汉字宽度及数字无千分节问题
func Sprintf(format string, args ...any) (d string) {
	StrPattern := regexp.MustCompile(`%(-)?(\d+)s`)
	IntPattern := regexp.MustCompile(`%(\d+)?,d`)
	FloatPattern := regexp.MustCompile(`%(\d+),\.(\d+)f`)
	ValuesPattern := regexp.MustCompile(`%(\d+)V`)
	Pattern := regexp.MustCompile(`%.*?[sdfvV%]`)
	i := 0
	replFunc := func(s string) (d string) {
		if s == "%%" {
			return "%"
		} else if k := StrPattern.FindStringSubmatch(s); k != nil {
			d = args[i].(string)
			l, _ := strconv.Atoi(k[2])
			if l-Wlen(d) > 0 {
				space := strings.Repeat(" ", l-Wlen(d))
				if k[1] == "-" {
					d = d + space
				} else {
					d = space + d
				}
			}
		} else if k := IntPattern.FindStringSubmatch(s); k != nil {
			l, _ := strconv.Atoi(k[1])
			d = string(formatInt(fmt.Sprintf("%d", args[i])))
			if l-len(d) > 0 {
				space := strings.Repeat(" ", l-len(d))
				d = space + d
			}

		} else if k := FloatPattern.FindStringSubmatch(s); k != nil {
			l, _ := strconv.Atoi(k[1])
			s, _ := strconv.Atoi(k[2])
			d = fmt.Sprintf(fmt.Sprintf("%%.%df", s), args[i])
			a := strings.Split(d, ".")
			a[0] = formatInt(a[0])
			d = strings.Join(a, ".")
			if l-len(d) > 0 {
				space := strings.Repeat(" ", l-len(d))
				d = space + d
			}
		} else if k := ValuesPattern.FindStringSubmatch(s); k != nil {
			count, _ := strconv.Atoi(k[1])
			d = strings.Repeat("?,", count-1)
			d += "?"
			return fmt.Sprintf("values(%s)", d)
		} else {
			d = fmt.Sprintf(s, args[i])
		}
		i++
		return
	}
	d = Pattern.ReplaceAllStringFunc(format, replFunc)
	return
}

// Printf 格式化打印
func Printf(format string, a ...any) {
	fmt.Print(Sprintf(format, a...))
}

// ChPrintln 打印通道中的数据
func ChPrintln[T any](ch <-chan []T) {
	for row := range ch {
		fmt.Println(Slice(row)...)
	}
}

// ChPrintf 格式化打印通道数据
func ChPrintf[T any](format string, ch <-chan []T, print_rows bool) {
	i := 0
	for row := range ch {
		Printf(format, Slice(row)...)
		i++
	}
	if print_rows {
		Printf("共 %,d 行数据\n", i)
	}
}

// PrintStruct 以 JSON 的格式打印结构化数据
func PrintStruct(v any) {
	if b, err := json.MarshalIndent(v, "", "    "); err == nil {
		fmt.Println(string(b))
	} else {
		CheckFatal(err)
	}
}

// GetMd5 获取指定字符的 md5 值
func GetMd5(strs ...string) string {
	b := []byte(strings.Join(strs, ""))
	return fmt.Sprintf("%x", md5.Sum(b))
}

var Sep = regexp.MustCompile(`[|,\t]| +`)

// split 分隔字符串
func Split(s string) []string {
	return Sep.Split(s, -1)
}
