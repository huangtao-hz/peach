// 项目：标准库，对 sqlite3 进行封装
// 模块：数据库
// 作者：黄涛
// 创建：2025-08-31
package sqlite

import (
	"database/sql"
	"peach/utils"
)

// sql 的查询接口，进行转接
type querier interface {
	Query(query string, args ...any) (rows *sql.Rows, err error)
}

// 查询结果
type Rows struct {
	*sql.Rows
}

// 执行查询
func Query(querier querier, query string, args ...any) (rows *Rows, err error) {
	var _rows *sql.Rows
	if _rows, err = querier.Query(query, args...); err != nil {
		return
	}
	rows = &Rows{_rows}
	return
}

// 提取所有结果数据，数据读取完成后自动关闭数据集
func (rows *Rows) FetchAll(ch chan<- []any) {
	defer rows.Close() // 关闭数据集
	defer close(ch)    // 关闭数据通道
	var addrs, values []any
	columns, err := rows.Columns()
	utils.CheckFatal(err)
	count := len(columns)
	values = make([]any, count)
	addrs = make([]any, count)
	for i := range count {
		addrs[i] = &values[i]
	}
	for rows.Next() {
		rows.Scan(addrs...)
		v := make([]any, count)
		copy(v, values)
		ch <- v
	}
}

// 逐行打印数据
func (r *Rows) Println() {
	ch := make(chan []any, 100)
	go r.FetchAll(ch)
	utils.ChPrintln(ch)
}

// 逐行格式化打印数据
func (r *Rows) Printf(format string, print_rows bool) {
	ch := make(chan []any, 100)
	go r.FetchAll(ch)
	utils.ChPrintf(format, ch, print_rows)
}
