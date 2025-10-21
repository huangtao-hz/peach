package main

import (
	"fmt"
	"peach/utils"
)

// Config 配置文件模型
type Config struct {
	Database string `toml:"database"`
	Home     string `toml:"home"`
}

// 配置文件，初始化的值
var config = Config{
	Database: "xgm2025-03",
	Home:     "~/Downloads",
}

// 模块初始化
func init() {
	utils.InitLog()
	if err := utils.GetConfig("xmjh", &config); err != nil {
		fmt.Println("导入配置文件失败，错误：", err)
	}
}
