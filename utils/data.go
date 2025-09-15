package utils

import (
	"context"
	"fmt"
)

// 数据通道，包含一个可取消的 Context 和 数据通道
type Data struct {
	context.Context
	Cancel context.CancelCauseFunc
	Data   chan []any
}

// 数据通道的构造函数
func NewData() *Data {
	ctx, cancel := context.WithCancelCause(context.Background())
	data := make(chan []any, 1000)
	return &Data{ctx, cancel, data}
}

// 打印数据通道中的数据
func (d *Data) Println() {
	for {
		select {
		case <-d.Done():
			return
		case row := <-d.Data:
			if row != nil {
				fmt.Println(row...)
			} else {
				return
			}
		}
	}
}

// 格式化打印数据
func (d *Data) Printf(format string) {
	for {
		select {
		case <-d.Done():
			fmt.Println("Error:", d.Err())
			return
		case row := <-d.Data:
			if row != nil {
				fmt.Print(Sprintf(format, row...))
			} else {
				return
			}
		}
	}
}
