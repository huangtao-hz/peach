package main

import (
	"fmt"
	"iter"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Item 数字表达式
type Item struct {
	level int // 0-数字，1-加减，2-乘除
	value int
	Repr  string
}

// NewItem 构造函数
func NewItem(value int) *Item {
	return &Item{value: value, Repr: strconv.Itoa(value)}
}

// String 显示表达式
func (i *Item) String() string {
	return fmt.Sprint(i.Repr, "=", i.value)
}

// Sorted 排序
func Sorted(a, b *Item) (*Item, *Item) {
	if a.value < b.value {
		a, b = b, a
	}
	return a, b
}

// Add 加法
func (i *Item) Add(other *Item) *Item {
	i, other = Sorted(i, other)
	return &Item{value: i.value + other.value, level: 1, Repr: fmt.Sprint(i.Repr, "+", other.Repr)}
}

// Sub 减法
func (i *Item) Sub(other *Item) *Item {
	i, other = Sorted(i, other)
	repr2 := other.Repr
	if other.level > 0 {
		repr2 = fmt.Sprint("(", repr2, ")")
	}
	return &Item{value: i.value - other.value, level: 1, Repr: fmt.Sprint(i.Repr, "-", repr2)}
}

// Multiply 乘法
func (i *Item) Multiply(other *Item) *Item {
	i, other = Sorted(i, other)
	repr1, repr2 := i.Repr, other.Repr
	if i.level == 1 {
		repr1 = fmt.Sprint("(", repr1, ")")
	}
	if other.level == 1 {
		repr2 = fmt.Sprint("(", repr2, ")")
	}
	return &Item{value: i.value * other.value, level: 2, Repr: fmt.Sprint(repr1, "*", repr2)}
}

// divide
func (i *Item) Divide(other *Item) *Item {
	i, other = Sorted(i, other)
	if other.value == 0 {
		return nil
	}
	repr1, repr2 := i.Repr, other.Repr
	if i.level == 1 {
		repr1 = fmt.Sprint("(", repr1, ")")
	}
	if other.level > 0 {
		repr2 = fmt.Sprint("(", repr2, ")")
	}
	if i.value%other.value == 0 {
		return &Item{value: i.value / other.value, level: 2, Repr: fmt.Sprint(repr1, "÷", repr2)}
	}
	return nil
}

func (i *Item) IterAll(other *Item) iter.Seq[*Item] {
	return func(yield func(*Item) bool) {
		values := []*Item{i.Add(other), i.Sub(other), i.Multiply(other)}
		if x := i.Divide(other); x != nil {
			values = append(values, x)
		}
		for _, x := range values {
			if !yield(x) {
				break
			}
		}
	}
}

// permute 生成数组序列
func permute(nums []int) [][]int {
	var res [][]int
	var path []int
	var used = make([]bool, len(nums))

	// 构建排列
	var backtrack func(int)
	backtrack = func(start int) {
		if start == len(nums) {
			res = append(res, append([]int{}, path...))
			return
		}
		for i := range nums {
			if used[i] {
				continue
			}
			path = append(path, nums[i])
			used[i] = true
			backtrack(start + 1)
			path = path[:len(path)-1]
			used[i] = false
		}
	}

	// 调用回溯函数
	backtrack(0)
	return res
}

// regular
func regular(s string) string {
	Jiafa := regexp.MustCompile(`\d+(\+\d+)+`)
	Chengfa := regexp.MustCompile(`\d+(\*\d+)+`)
	jffn := func(s string) string {
		var i []int
		for k := range strings.SplitSeq(s, "+") {
			x, _ := strconv.Atoi(k)
			i = append(i, x)
		}
		slices.Sort(i)
		var ss []string
		for _, k := range i {
			ss = append(ss, strconv.Itoa(k))
		}
		return strings.Join(ss, "+")
	}
	cffn := func(s string) string {
		var i []int
		for k := range strings.SplitSeq(s, "*") {
			x, _ := strconv.Atoi(k)
			i = append(i, x)
		}
		slices.Sort(i)
		var ss []string
		for _, k := range i {
			ss = append(ss, strconv.Itoa(k))
		}
		return strings.Join(ss, "*")
	}
	s = Jiafa.ReplaceAllStringFunc(s, jffn)
	s = Chengfa.ReplaceAllStringFunc(s, cffn)
	return s
}

// main 主函数
func main() {
	args := os.Args[1:]
	nums := make([]int, 4)
	result := make(map[string]bool)
	var (
		err error
	)
	if len(args) == 4 {
		for i, a := range args {
			if nums[i], err = strconv.Atoi(a); err != nil {
				return
			}
		}
	} else {
		return
	}
	for _, x := range permute(nums) {
		items := make([]*Item, 4)
		for i, k := range x {
			items[i] = NewItem(k)
		}
		for c := range items[0].IterAll(items[1]) {
			for d := range c.IterAll(items[2]) {
				for e := range d.IterAll(items[3]) {
					if e.value == 24 {
						result[regular(e.String())] = true
					}
				}
			}
			for d := range items[2].IterAll(items[3]) {
				for e := range c.IterAll(d) {
					if e.value == 24 {
						result[regular(e.String())] = true
					}
				}
			}
		}
	}

	for r := range result {
		fmt.Println(r)
	}
}
