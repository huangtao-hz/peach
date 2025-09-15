package main

import (
	"fmt"
	"peach/sqlite"
)

const create_sql = `
-- 项目计划表
create table if not exists xmjh(
	jym		text	primary key,	-- 交易码
	jymc	text,	-- 交易名称
	jyz 	text,	-- 交易组
	jyzm	text,	-- 交易组名
	yjcd	text,	-- 一级菜单
	ejcd	text,	-- 二级菜单
	ywl		int,	-- 业务量
	lx		text,	-- 类型
	bm 		text,	-- 业务部门
	zx		text,	-- 中心
	lxr		text,	-- 联系人
	fa 		text,	-- 方案
	xqwc	text,	-- 需求完成时间
	dqjd	text,	-- 当前进度
	bz		text,	-- 备注
	xjy		text	-- 对应新交易
);

-- 开发计划表
create table if not exists kfjh(
	jym		text	primary key, -- 交易码
	xqzt	text,	-- 需求状态
	kfzt	text,	-- 开发状态
	jhbb	text,	-- 计划版本
	kffzr	text,	-- 行方开发负责人
	kfzz	text,	-- 开发组长
	qdkf	text,	-- 前端开发
	hdkf	text,	-- 后端开发
	lckf	text,	-- 流程开发
	jccsks	text,	-- 开始集成测试
	jccsjs	text,	-- 结束即成测试
	yscsks	text,	-- 开始验收测试
	yscsjs	text 	-- 结束验收测试
);

-- 新旧交易对照表
create table if not exists jydzb(
	xjym	text,	-- 新交易码
	xjymc	text,	-- 新交易名称
	jym		text,	-- 老交易码
	jymc	text,	-- 老交易名称
	tcrq	text,	-- 投产日期
	zt 		text,	-- 状态
	bz		text,	-- 备注
	primary key(xjym,jym)
);
create index if not exists jydzb_tcrq on jydzb(tcrq);
`

func CreateDatabse(db *sqlite.DB) {
	fmt.Println("初始化数据库表")
	db.ExecScript(create_sql)
	sqlite.InitLoadFile(db)
	fmt.Println("初始化数据库成功！")
}
