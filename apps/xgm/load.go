package main

import (
	"fmt"
	"peach/sqlite"
	"peach/utils"
	"strings"
	"time"
)

// Update 更新计划表
func Update(db *sqlite.DB) (err error) {
	load_qxzb(db)
	path := utils.NewPath(config.Home).Find("*新柜面存量交易迁移*.xlsx")
	if path != nil {
		load_xmjh(db, path)
	} else {
		path = utils.NewPath(config.Home).Join(fmt.Sprintf("附件1：新柜面存量交易迁移计划%s.xlsx", utils.Today().Format("%Y%M%D")))
	}
	//load_kfjh(db) // 导入科技管理部编制的开发计划表
	fmt.Print("根据投产时间更新验收完成时间:")
	r, _ := db.Exec(`update bbap set wcys=date(tcrq,"weekday 5","-7 days") where wcys=""`)
	if count, err := r.RowsAffected(); err == nil {
		fmt.Println(count, "条数据被更新")
	}
	Update_ytc(db)
	fmt.Print("根据验收明细表更新开发状态:")
	db.ExecuteFs(queryFS, "query/update_kfjihua.sql")
	fmt.Print("根据计划版本更新开发计划时间：")
	db.ExecuteFs(queryFS, "query/update_kfjhsj.sql")
	fmt.Print("根据验收条目更新完成状态：")
	db.ExecuteFs(queryFS, "query/update_xmjh.sql")
	fmt.Print("根据新旧交易对照表更新对应新交易：")
	db.ExecuteFs(queryFS, "query/update_xmjh_xjy.sql")
	Export(db, path)
	return
}

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
	var path *utils.Path
	if path = utils.NewPath(config.Home).Find("新柜面简报*.zip"); path == nil {
		return fmt.Errorf("未找到 新柜面简报*.zip 文件")
	}
	fmt.Println("处理文件：", path.Name())
	for name, file := range path.IterZip() {
		if strings.Contains(name, "新柜面存量交易迁移计划") {
			load_xmjh(db, file)
		} else if strings.Contains(name, "数智综合运营系统问题跟踪表") {
			load_wtgzb(db, file)
		} else if strings.Contains(name, "版本条目明细") {
			load_bbmx(db, file)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}
