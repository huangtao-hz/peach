package data

import (
	"fmt"
	"peach/utils"
	"slices"
)

// ConvertFunc 数据转换函数类型
type ConvertFunc func([]string) ([]string, error)

// FixedColumn 固定列数
func FixedColumn(count int) ConvertFunc {
	return func(source []string) (dest []string, err error) {
		if len(source) >= count {
			dest = source[:count]
		} else if len(source) < count {
			dest = append(source, slices.Repeat([]string{""}, count-len(source))...)
		}
		return
	}
}

// Include 包含指定列
func Include(columns ...int) ConvertFunc {
	dest_length := len(columns)
	return func(source []string) (dest []string, err error) {
		dest = make([]string, dest_length)
		source_length := len(source)
		if source_length < dest_length {
			err = fmt.Errorf("源切片长度不足")
			return
		}
		for i, k := range columns {
			if k < 0 {
				k += source_length
			} else if k < source_length {
				dest[i] = source[k]
			}
		}
		fmt.Println(len(dest), dest)
		return
	}
}

// Exclude 剔除指定列
func Exclude(columns ...int) ConvertFunc {
	return func(source []string) (dest []string, err error) {
		dest = make([]string, 0)
		new_columns := make([]int, len(columns))
		source_length := len(source)
		for i, k := range columns {
			if k < 0 {
				new_columns[i] = k + source_length
			}
		}
		for i, k := range source {
			if slices.Index(new_columns, i) < 0 {
				dest = append(dest, k)
			}
		}
		return
	}
}

// Convert 转换函数
func Convert(covertvertFuncs ...ConvertFunc) func([]string) ([]any, error) {
	return func(src []string) (dest []any, err error) {
		for _, convertFunc := range covertvertFuncs {
			src, err = convertFunc(src)
			if src == nil || err != nil {
				return
			}
		}
		dest = utils.Slice(src)
		return
	}
}
