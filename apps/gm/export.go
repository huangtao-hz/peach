package main

import (
	"fmt"
	"peach/data"
	"peach/excel"
	"peach/sqlite"
	"peach/utils"
)

const (
	KfjhHeader = "交易码,交易名称,交易组,交易组名,一级菜单,二级菜单,近一年交易量,类型,部门,中心,联系人,方案,需求状态,开发状态,计划版本,开发负责人,开发组长,前端开发,后端开发,流程开发,集成测试开始,集成测试结束,验收测试开始,验收测试结束"
	KfjhQuery  = `
select a.jym,a.jymc,a.jyz,a.jyzm,a.yjcd,a.ejcd,a.bs,a.lx,a.ywbm,a.zx,a.lxr,a.fa,
b.xqzt,b.kfzt,b.jhbb,b.kjfzr,b.kfzz,b.qdkf,b.hdkf,b.lckf,b.jcks,b.jcjs,b.ysks,b.ysjs
from xmjh a
left join kfjh b on a.jym=b.jym
where b.jym is not null
order by b.jhbb,a.jym
`
	BufferSize = 1024
)

var KfjhWidth = map[string]float64{
	"A,C":       6.83,
	"B":         42,
	"D":         21,
	"E,S,T":     15,
	"F":         31,
	"G,I,J,P,R": 14,
	"H,K":       9,
	"L:O,Q":     11,
	"U:X":       16,
}
var KfjhStyle = map[string]string{
	"A:F,H:X": "Normal-NoWrap",
	"G":       "Number",
}

func export_xjdzb(db *sqlite.DB, book *excel.Writer) (err error) {
	sheet := book.GetSheet("投产交易一览表")
	header := "交易码,交易码称,原交易码,原交易码称,投产日期,状态,备注"
	sheet.SetWidth(map[string]float64{
		"A,C": 9,
		"B,D": 42,
		"E":   11,
		"F":   12,
		"G":   17,
	})
	sheet.SetColStyle(map[string]string{
		"A:G": "Normal-NoWrap",
	})
	query := "select * from xjdz order by jym,yjym"
	rows, err := db.Query(query)
	if err != nil {
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.WriteTable("A1", header, ch)
	return
}

func export_kfjh(db *sqlite.DB, book *excel.Writer) (err error) {
	sheet := book.GetSheet("开发计划")
	sheet.SetWidth(KfjhWidth)
	sheet.SetColStyle(KfjhStyle)
	rows, err := db.Query(KfjhQuery)
	if err != nil {
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.WriteTable("A1", KfjhHeader, ch)
	return
}

func export_xmjh(db *sqlite.DB, book *excel.Writer) (err error) {
	header := "交易码,交易名称,交易组,交易组名,一级菜单,二级菜单,近一年交易量,类型,部门,中心,联系人,方案,计划需求完成时间,当前进度,备注,新交易"
	querys := map[string]string{
		"计划表": "select * from xmjh where sfwc is null or not sfwc like '5%' order by jym",
		"完成表": "select * from xmjh where sfwc like '5%' order by jym",
		"全量表": "select * from xmjh order by jym",
	}
	for name, query := range querys {
		sheet := book.GetSheet(name)
		sheet.SetWidth(map[string]float64{
			"A,C":       6.83,
			"B":         42,
			"D":         21,
			"E":         15,
			"F":         31,
			"G,I,J,P,N": 14,
			"H,K,M":     9,
			"O":         24,
			"L":         11,
			"Q":         35,
		})
		sheet.SetColStyle(map[string]string{
			"A:F,H:P": "Normal-NoWrap",
			"G":       "Number",
		})
		sheet.SetColVisible("Q", false)
		if rows, err := db.Query(query); err == nil {
			inch := make(chan []any, BufferSize)
			outch := make(chan []any, BufferSize)
			go rows.FetchAll(inch)
			go func() {
				defer close(outch)
				var hasher = data.Hashier(-9, -8, -7, -6, -5, -4, -3, -2, -1)
				for row := range inch {
					dest := make([]string, len(row))
					for i, k := range row {
						dest[i] = fmt.Sprintf("%v", k)
					}
					dest, err = hasher(dest)
					outch <- utils.Slice(dest)
				}
			}()
			sheet.WriteTable("A1", header, outch)
		} else {
			return err
		}
	}
	return
}

const (
	tongji_sql = `
select lxr,zx,
sum(iif((sfwc is null or sfwc ='0-尚未开始' or sfwc='' ),1,0)),       -- 未开始
sum(iif(sfwc in('1-已编写初稿','2-已提交需求/确认需规'),1,0)),       -- 已完成需求
sum(iif(sfwc in('3-已完成开发','4-已完成验收测试'),1,0)),       -- 开发中
sum(iif(sfwc = '5-已投产' ,1,0)) ,       -- 已投产
count(jym) as zs,         -- 总数
sum(iif(fa <> '1-下架交易' and sfwc = '5-已投产',1,0))*1.0/count(jym)
from xmjh
where ywbm='运营管理部' and fa not in('1-下架交易','5-移出柜面系统')
group by zx,lxr
order by zs desc
`

	tongji_gzx_sql = `
select zx,
sum(iif((sfwc is null or sfwc ='0-尚未开始' or sfwc='' ),1,0)),       -- 未开始
sum(iif(sfwc in('1-已编写初稿','2-已提交需求/确认需规'),1,0)),       -- 已完成需求
sum(iif(sfwc in('3-已完成开发','4-已完成验收测试'),1,0)),       -- 开发中
sum(iif(sfwc = '5-已投产' ,1,0)) ,       -- 已完成需求
count(jym) as zs,         -- 总数
sum(iif(fa <> '1-下架交易' and sfwc = '5-已投产',1,0))*1.0/count(jym)
from xmjh
where ywbm='运营管理部' and fa not in('1-下架交易','5-移出柜面系统')
group by zx
order by zs desc
`

	tongji_gbm_sql = `
select lx,
sum(iif(fa <> '1-下架交易' and(sfwc is null or sfwc ='0-尚未开始' or sfwc=''),1,0)),       -- 未开始
sum(iif(fa <> '1-下架交易' and sfwc in('1-已编写初稿','2-已提交需求/确认需规'),1,0)),       -- 已完成需求
sum(iif(fa <> '1-下架交易' and sfwc in('3-已完成开发','4-已完成验收测试'),1,0)),       -- 开发中
sum(iif(fa <> '1-下架交易' and sfwc = '5-已投产',1,0)),       -- 已完成需求
count(jym) as zs,         -- 总数
sum(iif(fa <> '1-下架交易' and sfwc = '5-已投产',1,0))*1.0/count(jym)
from xmjh
where fa not in('1-下架交易','5-移出柜面系统')
group by lx
order by zs desc
`
)

func export_tjb(db *sqlite.DB, book *excel.Writer) (err error) {
	sheet := book.GetSheet(0)
	sheet.Rename("统计表")
	header := "类型,未提交需求,已完成需求,开发中,已投产,总数,投产完成率"
	sheet.SetWidth(map[string]float64{
		"A":   12,
		"B":   20,
		"C:H": 10,
	})
	sheet.SetColStyle(map[string]string{
		"A:B": "Normal-NoWrap",
		"C:G": "Number",
		"H":   "Percent",
	})

	rows, err := db.Query(tongji_gbm_sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch := make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.WriteTable("B1", header, ch)

	header = "中心,未提交需求,已完成需求,开发中,已投产,总数,投产完成率"
	rows, err = db.Query(tongji_gzx_sql)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch = make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.WriteTable("B8", header, ch)

	header = "联系人,中心,未提交需求,已完成需求,开发中,已投产,总数,投产完成率"
	rows, err = db.Query(tongji_sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch = make(chan []any, BufferSize)
	go rows.FetchAll(ch)
	sheet.WriteTable("A19", header, ch)

	return
}
func Export(db *sqlite.DB, path *utils.Path) (err error) {
	fmt.Println("更新文件：", path)
	book := excel.NewWriter()
	export_tjb(db, book)
	export_kfjh(db, book)
	export_xmjh(db, book)
	export_xjdzb(db, book)
	book.SaveAs(path.String())
	fmt.Println("更新文件完成！")
	return
}
