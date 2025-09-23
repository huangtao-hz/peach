package main

import (
	_ "embed"
	"fmt"
	"peach/sqlite"
)

//go:embed xmjh.sql
var create_sql string

// CreateDatabse 创建数据库表
func CreateDatabse(db *sqlite.DB) {
	fmt.Println("初始化数据库表")
	db.ExecScript(create_sql)
	sqlite.InitLoadFile(db)
	fmt.Println("初始化数据库成功！")
}
