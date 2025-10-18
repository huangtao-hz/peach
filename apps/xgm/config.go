package main

import (
	"fmt"
	"peach/utils"
)

type Config struct {
	Database string `toml:"database"`
	Home     string `toml:"home"`
}

var config = Config{
	Database: "xgm2025-03",
	Home:     "~/Downloads",
}

func init() {
	utils.InitLog()
	if err := utils.GetConfig("xmjh", &config); err != nil {
		fmt.Println("导入配置文件失败，错误：", err)
	}
}
