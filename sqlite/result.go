// 项目：标准库，对 sqlite3 进行封装
// 模块：数据库
// 作者：黄涛
// 创建：2025-08-31
package sqlite

type result struct {
	rowsAffected int64
	lastInsertId int64
	err          error
}

// 返回最后插入的 Id
func (r *result) LastInsertId() (int64, error) {
	return r.lastInsertId, r.err
}

// 返回影响的行数
func (r *result) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}
