// 项目：标准库，对 sqlite3 进行封装
// 模块：数据库
// 作者：黄涛
// 创建：2025-08-31

package sqlite

import (
	"database/sql"
	"fmt"
	"peach/utils"
	"strings"

	//_ "modernc.org/sqlite"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/ncruces/go-sqlite3/driver"
	//_ "github.com/ncruces/go-sqlite3/embed"
)

// 数据库连接
type DB struct {
	*sql.DB
}

// 打开数据库连接
func Open(database string) (db *DB, err error) {
	if (database != ":memory:") && (utils.NewPath(database).Dir() == ".") {
		dataHome := utils.Home.Join(".data")
		dataHome.Ensure() // 目录不存在则自动创建
		database = (dataHome.Join(database).WithExt(".db")).String()
	}
	if _db, err := sql.Open("sqlite3", database); err == nil {
		db = &DB{_db}
	}
	return
}

// 执行查询
func (db *DB) Query(query string, args ...any) (*Rows, error) {
	return Query(db.DB, query, args...)
}

// 执行 DDL 语句，一般用来生成库表
func (db *DB) ExecScript(query string) {
	_, err := db.Exec(query)
	utils.CheckFatal(err)
}

// 开启事物
func (db *DB) Begin() (tx *Tx, err error) {
	if tx_, err := db.DB.Begin(); err == nil {
		tx = &Tx{tx_}
	}
	return
}

// 执行 sql 语句
type ExecFunc func(*Tx) error

// 执行事务
func (db *DB) ExecTx(execFuncs ...ExecFunc) (err error) {
	tx, err := db.Begin() // 开启事务
	if err != nil {
		return
	}
	defer tx.Rollback() // 执行失败则回滚数据
	for _, execFunc := range execFuncs {
		if err = execFunc(tx); err != nil {
			return
		}
	}
	return tx.Commit() // 无异常则提交数据库
}

// 执行查询，并格式化打印结果
func (db *DB) Printf(query string, format string, head string, print_rows bool, args ...any) {
	r, err := db.Query(query, args...)
	utils.CheckFatal(err)
	if head != "" {
		fmt.Println(head)
	}
	r.Printf(format, print_rows)
}

// 执行查询，并格式化打印结果
func (db *DB) Println(query string, args ...any) {
	r, err := db.Query(query, args...)
	utils.CheckFatal(err)
	r.Println()
}

// 打印一行数据，采用 Key-Value 的格式输出
func (db *DB) PrintRow(query string, header string, args ...any) (err error) {
	var width int
	headers := strings.Split(header, ",")
	count := len(headers)
	values := make([]any, count)
	addrs := make([]any, count)
	for i := range headers {
		addrs[i] = &values[i]
		if l := utils.Wlen(headers[i]); l > width {
			width = l
		}
	}
	format := fmt.Sprintf("%%%ds  %%s", width)
	err = db.QueryRow(query, args...).Scan(addrs...)
	if err != nil {
		return
	}
	for i, header := range headers {
		fmt.Println(utils.Sprintf(format, header, fmt.Sprint(values[i])))
	}
	return
}
