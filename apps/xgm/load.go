package main

import (
	"fmt"
	"peach/sqlite"
	"peach/utils"
	"strings"
	"time"
)

// Load 导入数据文件
func Load(db *sqlite.DB) (err error) {
	Home := utils.NewPath(config.Home)
	if path := Home.Find("*新柜面存量交易迁移*.xlsx"); path != nil {
		load_xmjh(db, path)
	}
	if path := Home.Find("*数智综合运营系统问题跟踪表*.xlsx"); path != nil {
		load_wtgzb(db, path)
	}
	return nil
}

// Restore 从备份文件中恢复数据
func Restore(db *sqlite.DB) (err error) {
	defer utils.TimeIt(time.Now())
	home := utils.NewPath(config.Home)
	var path *utils.Path
	if path = home.Find("新柜面简报*.zip"); path == nil {
		return fmt.Errorf("未找到 新柜面简报*.zip 文件")
	}
	fmt.Println("处理文件：", path.Name())
	for name, file := range path.IterZip() {
		if strings.Contains(name, "新柜面存量交易迁移计划") {
			load_xmjh(db, file)
			f := home.Join(file.FileInfo().Name())
			export_xmjh(db, f)
		} else if strings.Contains(name, "数智综合运营系统问题跟踪表") {
			load_wtgzb(db, file)
		} else if strings.Contains(name, "版本条目明细") {
			load_bbmx(db, file)
			f := home.Join(file.FileInfo().Name())
			export_bbmx(db, f)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}
