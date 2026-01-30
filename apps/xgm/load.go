package main

import (
	"fmt"
	"peach/utils"
	"strings"
	"time"
)

// Restore 从备份文件中恢复数据
func (c *Client) Restore() (err error) {
	defer utils.TimeIt(time.Now())
	home := utils.NewPath(c.Home)
	var path *utils.Path
	if path = home.Find("新柜面简报*.zip"); path == nil {
		return fmt.Errorf("未找到 新柜面简报*.zip 文件")
	}
	fmt.Println("处理文件：", path.Name())
	for name, file := range path.IterZip() {
		if strings.Contains(name, "新柜面存量交易迁移计划") {
			c.load_xmjh(file)
			f := home.Join(file.FileInfo().Name())
			c.export_xmjh(f)
		} else if strings.Contains(name, "数智综合运营系统问题跟踪表") {
			c.load_wtgzb(file)
		} else if strings.Contains(name, "版本条目明细") {
			c.load_bbmx(file)
			f := home.Join(file.FileInfo().Name())
			c.export_bbmx(f)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}
