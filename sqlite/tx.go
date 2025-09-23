// 项目：标准库，对 sqlite3 进行封装
// 模块：数据库
// 作者：黄涛
// 创建：2025-08-31
package sqlite

import (
	"database/sql"
	"fmt"
	"peach/data"
)

// Tx 数据库事务
type Tx struct {
	*sql.Tx
}

// ExecMany 批量执行
func (t *Tx) ExecMany(query string, data *data.Data) (r sql.Result, err error) {
	stmt, err := t.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var (
		_result sql.Result
		c       int64
		count   int64 = 0
		lastId  int64
	)
	for {
		select {
		case <-data.Done():
			err = data.Err()
			return
		case row := <-data.Data:
			if row == nil {
				lastId, err = _result.LastInsertId()
				r = &result{count, lastId, err}
				return
			} else if _result, err = stmt.ExecContext(data.Context, row...); err == nil {
				c, _ = _result.RowsAffected()
				count += c
			} else {
				return
			}
		}
	}
}

// Query 查询多行数据
func (t *Tx) Query(query string, args ...any) (rows *Rows, err error) {
	return Query(t.Tx, query, args...)
}

// GetColumnCount 查询指定数据库表的字段数量
func (t *Tx) GetColumnCount(tablename string) (count int, err error) {
	var rows *Rows
	rows, err = t.Query(fmt.Sprintf("select * from %s where 0", tablename))
	if err != nil {
		return
	}
	columns, err := rows.Columns()
	if err == nil {
		count = len(columns)
	}
	return
}
