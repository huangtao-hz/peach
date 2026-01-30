package main

import (
	"fmt"
	"peach/utils"
)

// conv_qxzb 缺陷指标转换程序
func conv_qxzb(src []string) (dest []string, err error) {
	if src[0] != "" {
		dest = src
	}
	return
}

// load_qxzb 导入缺陷指标
func (c *Client) load_qxzb() {
	if path := utils.NewPath(c.Home).Find("*缺陷指标详情*.xlsx"); path != nil {
		fmt.Println("处理文件：", path.Base())
		c.LoadExcelFile(loaderFS, "loader/qx_qxzb.toml", path, "", conv_qxzb)
	} else {
		fmt.Println("Error: 未发现 缺陷指标详情 文件")
	}
}
