// 项目：标准库，对 sqlite3 进行封装
// 模块：数据库
// 作者：黄涛
// 创建：2025-08-31
package sqlite

import (
	"database/sql"
)

// 回调函数
type CallbackFunc func(sql.Result, error) error

// 执行单条 sql 语句
func ExecOne(query string, callback CallbackFunc, args ...any) ExecFunc {
	return func(tx *Tx) (err error) {
		r, err := tx.Exec(query, args...)
		if callback != nil {
			return callback(r, err)
		}
		return
	}
}

// 执行多条语句
func ExecMany(query string, callback CallbackFunc, ch <-chan []any) ExecFunc {
	return func(tx *Tx) (err error) {
		r, err := tx.ExecMany(query, ch)
		if callback != nil {
			return callback(r, err)
		}
		return
	}
}
