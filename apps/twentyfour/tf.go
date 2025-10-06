package main

import (
	"fmt"
	"strconv"
)

type Operator int

const (
	Num     Operator = iota // 数字
	Sum                     // 求和
	Product                 // 求积
)

// Number 定义项目的机构体
type Number struct {
	value    int
	operator Operator // 0-number, 1-sum, 2-product
	reversed bool
	childs   []Number
}

// New 构造函数
func (n *Number) New(value int) *Number {
	return &Number{value: value, operator: Num}
}

func (n *Number) String() string {
	switch n.operator {
	case Num:
		return strconv.Itoa(n.value)
	case Sum:
		s := ""
		for _, k := range n.childs {
			st := k.String()
			if k.reversed {
				st = "-" + st
			} else {
				st = "+" + st
			}
			s += st[1:]
		}
	case Product:
		s := ""
		for _, k := range n.childs {
			st := k.String()
			if k.reversed {
				st = "÷" + st
			} else {
				st = "×" + st
			}
			s += st[1:]
		}
		return s
	}
	return ""
}

func main() {
	fmt.Println("Hello world.")
}
