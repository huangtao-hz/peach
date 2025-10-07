// 项目：标准库，对 sqlite3 进行封装
// 模块：数据库
// 作者：黄涛
// 创建：2025-08-31

package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"peach/data"
	"peach/utils"
)

const loadFileSQL = `create table if not exists loadfile(
name	text 	primary key,  -- 类型
path	text,	              -- 文件名
mtime	text,                 -- 文件修改时间
ver		text		          -- 文件版本
)`

// InitLoadFile 创建重复导入检查的数据库表
func InitLoadFile(db *DB) {
	db.Exec(loadFileSQL)
}

// Loader 数据装入类
type Loader struct {
	tablename, LoadSQL string          // 表名，导入的SQL语句
	fileinfo           os.FileInfo     // 文件信息
	Ver                string          // 导入数据版本
	reader             data.DataReader // 数据读取程序
	Check              bool            // 是否需要检查文件重复导入
	Clear              bool            // 是否清理数据库，默认为是
	ClearSQL           string          // 自定义清理语句
	Fields             string          // 插入字段清单，用半角逗号分隔
	FieldCount         int             // 字段总数，如有 Fields，优先使用 Fields
	Method             string          // 导入方法，默认为 insert，可以为：replace
	db                 *DB
}

// GetLoadSQL 生成导入的 sql 语句
func (l *Loader) GetLoadSQL(tx *Tx) string {
	if l.LoadSQL != "" {
		return l.LoadSQL
	}
	var (
		fields  []string                // 字段列表
		count   int      = l.FieldCount // 字段个数
		builder strings.Builder
	)
	if l.Fields != "" {
		fields = strings.Split(l.Fields, ",")
		count = len(fields)
	} else if count == 0 {
		// 从数据库中获取字段数
		count, _ = tx.GetColumnCount(l.tablename)
	}
	builder.WriteString(l.Method)
	builder.WriteString(" into ")
	builder.WriteString(l.tablename)
	if fields != nil {
		builder.WriteString(fmt.Sprintf("(%s)", strings.Join(fields, ",")))
	}
	builder.WriteString(" values(")
	builder.WriteString(strings.Join(slices.Repeat([]string{"?"}, count), ","))
	builder.WriteString(")")
	return builder.String()
}

// DoCheck 重复导入检查
func (l *Loader) DoCheck(tx *Tx) (err error) {
	if l.Check {
		var count int
		filename := l.fileinfo.Name()
		mtime := l.fileinfo.ModTime().Format("2006-01-02 15:04:05")
		const (
			checkSQL = "select count(name) from loadfile where name=? and path=? and mtime>=datetime(?)"
			doneSQL  = "insert or replace into loadfile values(?,?,datetime(?),?)"
		)
		err = tx.QueryRow(checkSQL, l.tablename, filename, mtime).Scan(&count)
		if err != nil {
			return
		}
		if count > 0 {
			return fmt.Errorf("文件 %s 已导入", filename)
		}
		_, err = tx.Exec(doneSQL, l.tablename, filename, mtime, l.Ver)
	}
	return
}

// DoTest 执行测试，完整跑通流程，但是不在数据库中执行
func (l *Loader) DoTest(tx *Tx) (err error) {
	fmt.Println("导入SQL：", l.GetLoadSQL(tx))
	return errors.New("测试成功")
}

// Test 测试导入数据
func (l *Loader) Test() {
	if err := l.db.ExecTx(l.DoCheck, l.DoClear, l.DoLoad, l.DoTest); err != nil {
		fmt.Println(err)
	}
}

// DoClear 执行清理
func (l *Loader) DoClear(tx *Tx) (err error) {
	if l.Clear && strings.Contains(l.Method, "insert") {
		if l.ClearSQL == "" {
			l.ClearSQL = fmt.Sprintf("delete from %s", l.tablename)
		}
		_, err = tx.Exec(l.ClearSQL)
	}
	return
}

// Success 导入成功后的提示语
func (l *Loader) Success(r sql.Result, err error) error {
	c, _ := r.RowsAffected()
	utils.Printf("文件 %s 导入成功, 共 %,d 行数据。\n", l.fileinfo.Name(), c)
	return nil
}

// DoLoad 执行导入流程
func (l *Loader) DoLoad(tx *Tx) (err error) {
	d := data.NewData()
	go l.reader.Read(d)
	return ExecMany(l.GetLoadSQL(tx), l.Success, d)(tx)
}

// Load 导入数据
func (l *Loader) Load() (err error) {
	if err = l.db.ExecTx(
		l.DoCheck, // 重复导入检查
		l.DoClear, // 清理历史数据
		l.DoLoad,  // 执行导入数据
	); err != nil {
		fmt.Println(err)
	}
	return
}
