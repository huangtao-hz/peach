select lx,
sum(iif(fa <> '1-下架交易' and(sfwc is null or sfwc ='0-尚未开始' or sfwc=''),1,0)),       -- 未开始
sum(iif(fa <> '1-下架交易' and sfwc in('1-已编写初稿','2-已提交需求/确认需规'),1,0)),       -- 已完成需求
sum(iif(fa <> '1-下架交易' and sfwc in('3-已完成开发','4-已完成验收测试'),1,0)),       -- 开发中
sum(iif(fa <> '1-下架交易' and sfwc = '5-已投产',1,0)),       -- 已完成需求
count(jym) as zs         -- 总数
from xmjh
where fa not in('1-下架交易','5-移出柜面系统')
group by lx
order by zs desc
