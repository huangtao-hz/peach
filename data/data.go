package data

import (
	"context"
	"fmt"
	"peach/utils"
)

const BuffSize = 1024

// DataReader 数据读取类
type DataReader interface {
	Read(*Data)
}

// 数据通道，包含一个可取消的 Context 和 数据通道
type Data struct {
	context.Context
	Cancel context.CancelCauseFunc
	Data   chan []any
}

// NewData 数据通道的构造函数
func NewData() *Data {
	ctx, cancel := context.WithCancelCause(context.Background())
	data := make(chan []any, BuffSize)
	return &Data{ctx, cancel, data}
}

// Println 打印数据通道中的数据
func (d *Data) Println() {
	for i := 0; ; i++ {
		select {
		case <-d.Done():
			return
		case row := <-d.Data:
			if row != nil {
				fmt.Println(i, len(row), row)
			} else {
				return
			}
		}
	}
}

// Printf 格式化打印数据
func (d *Data) Printf(format string) {
	for {
		select {
		case <-d.Done():
			fmt.Println("Error:", d.Err())
			return
		case row := <-d.Data:
			if row != nil {
				fmt.Print(utils.Sprintf(format, row...))
			} else {
				return
			}
		}
	}
}

// Cause 返回失败原因
func (d *Data) Err() error {
	return context.Cause(d.Context)
}
