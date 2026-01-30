package main

import (
	_ "embed"
	"flag"
	"fmt"
	"peach/sqlite"
	"peach/utils"
)

var Version = "1.0.10"

func CreateDatabse(db *sqlite.DB) {
	db.ExecFs(queryFS, "query/db.sql")
	sqlite.InitLoadFile(db)
}

// Client 定义客户端
type Client struct {
	*sqlite.DB
	database string `toml:"database"`
	Home     string `toml:"home"`
}

// Open 打开客户端
func Open() (client *Client, err error) {
	utils.InitLog()
	client = &Client{database: "xgm2025-03", Home: "~/Downloads"}
	if err = utils.GetConfig("xmjh", &client); err != nil {
		return
	}
	if client.DB, err = sqlite.Open(client.database); err != nil {
		return
	}
	sqlite.InitLoadFile(client.DB)         // 初始化导入文件数据库表
	client.ExecFs(queryFS, "query/db.sql") //初始化数据库表
	return
}

// show_version 显示程序版本
func (c *Client) show_version(string) (err error) {
	fmt.Println("版本：", Version)
	return
}

// load 导入数据
func (c *Client) load(string) (err error) {
	Home := utils.NewPath(c.Home)
	if path := Home.Find("*新柜面存量交易迁移*.xlsx"); path != nil {
		c.load_xmjh(path)
	}
	if path := Home.Find("*数智综合运营系统问题跟踪表*.xlsx"); path != nil {
		c.load_wtgzb(path)
	}
	return
}

// query_sql 执行查询语句
func (c *Client) query_sql(query *string) {
	c.Println(*query)

}
func (c *Client) update(string) (err error) {
	c.load_qxzb()
	c.update_bbmx()
	err = c.update_xmjh()
	Home := utils.NewPath(c.Home)
	if path := Home.Find("*数智综合运营系统问题跟踪表*.xlsx"); path != nil {
		c.load_wtgzb(path)
	}
	return
}

// Run 运行主程序
func (c *Client) Run() {
	flag.BoolFunc("version", "显示程序版本", c.show_version)
	flag.BoolFunc("update", "更新文件", c.update)
	flag.BoolFunc("touchan", "显示投产情况", c.show_touchan)
	flag.BoolFunc("jihua", "显示投产计划", c.kaifajihua)
	query_sql := flag.String("query", "", "执行查询")
	jhbb := flag.String("jhbb", "", "查询计划版本")
	wenti := flag.String("wenti", "", "统计上报问题，取值：本月、上月、上周、本周")
	flag.Parse()
	if *query_sql != "" {
		c.query_sql(query_sql)
	}
	if *jhbb != "" {
		c.show_jhbb(*jhbb)
	}
	if *wenti != "" {
		report_wenti(c.DB, *wenti)
	}
	for _, jym := range flag.Args() {
		if utils.FullMatch(`\d{5}`, jym) {
			c.show_new_jy(jym)
		} else if utils.FullMatch(`\d{4}`, jym) {
			c.show_old_jy(jym)
		}
	}
}

// main 主程序入口
func main() {
	defer utils.Recover()
	client, err := Open()
	utils.CheckFatal(err)
	defer client.Close()
	client.Run()
}

/*
func main_bak() {
	defer utils.Recover()
	db, err := sqlite.Open(config.Database)
	utils.CheckFatal(err)
	defer db.Close()
	load := flag.Bool("load", false, "导入数据")
	query_sql := flag.String("query", "", "执行查询")
	jhbb := flag.String("jhbb", "", "查询计划版本")
	restore := flag.Bool("restore", false, "导入数据")
	touchan := flag.Bool("touchan", false, "导入数据")
	update := flag.Bool("update", false, "更新计划表")
	jihua := flag.Bool("jihua", false, "投产交易清单")
	version := flag.Bool("version", false, "查阅程序版本")
	wenti := flag.String("wenti", "", "统计上报问题，取值：本月、上月、上周、本周")

	flag.Parse()
	CreateDatabse(db)
	if *version {
		fmt.Println("版本：", Version)
	}
	if *load {
		//err = load_jyjh(db)
		err = Load(db)
	}
	if *query_sql != "" {
		db.Println(*query_sql)
	}
	if *jhbb != "" {
		show_jhbb(db, *jhbb)
	}
	if *restore {
		Restore(db)
	}
	if *touchan {
		show_touchan(db)
	}
	if *jihua {
		kaifajihua(db)
	}
	if *wenti != "" {
		report_wenti(db, *wenti)
	}
	if len(flag.Args()) > 0 {
		PrintVersion(db)
	}
	for _, jym := range flag.Args() {
		if utils.FullMatch(`\d{5}`, jym) {
			show_new_jy(db, jym)
		} else if utils.FullMatch(`\d{4}`, jym) {
			show_old_jy(db, jym)
		}
	}
	utils.CheckErr(err)
}

*/
