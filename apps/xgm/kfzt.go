package main

import (
	"database/sql"
	"fmt"
	"peach/data"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
	"strings"
)

func convdate(d *string) {
	if strings.Contains(*d, "/") {
		*d = strings.ReplaceAll(*d, "/", "-")
	}
	if len(*d) > 0 && len(*d) < 10 {
		*d = ""
	}
}

// conv_kfzt 开发状态转换
func conv_kfazt(src []string) (dest []string, err error) {
	if src[10] == "" {
		return
	} else if len(src[10]) < 4 {
		src[10] = fmt.Sprintf("%04s", src[10])
	}
	for i := 1; i < 6; i++ {
		src[i] = strings.TrimSpace(src[i])
	}
	for i := 6; i < 10; i++ {
		convdate(&src[i])
	}
	dest = src
	return
}

// update_kfzt 更新开发状态
func update_kfzt(db *sqlite.DB) (err error) {
	path := utils.NewPath("~/Downloads").Find("*开发计划*.xlsx")
	if path == nil {
		return fmt.Errorf("未找到 开发计划 文件")
	}
	fmt.Println("处理文件：", path.Base())
	f, err := excel.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	r, err := f.NewReader("柜面核心类交易开发计划", "U,AE,AF,AP,AW,BD,AG,AH,BM,BS,B", 1, conv_kfazt)
	if err != nil {
		return
	}
	d := data.NewData()

	query := "update kfjh set kfzt=?,kjfzr=?,kfzz=?,qdkf=?,hdkf=?,lckf=?,jcks=?,jcjs=?,ysks=?,ysjs=? where jym=?"
	var callback = func(r sql.Result, err error) error {
		if err == nil {
			rows, _ := r.RowsAffected()
			utils.Printf("Total: %,d rows affected.\n", rows)
		} else {
			fmt.Println(err)
		}
		return err
	}

	go r.Read(d)
	err = db.ExecTx(sqlite.ExecMany(query, callback, d))
	return
}
