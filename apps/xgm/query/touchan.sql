-- 投产交易清单，列表显示投产交易的统计数据
select tcrq,count(distinct jym),sum(iif(yjym<>'',1,0)),sum(iif(yjym='',1,0)),
sum(iif(yjym<>'',1,0))*100.0/(select count(jym)from xmjh where fa not in ('1-下架交易','5-移出柜面系统'))
from xjdz
where tcrq<=date('now')
group by tcrq
union
-- 显示合计数据
select '合计',count(distinct jym),sum(iif(yjym<>'',1,0)),sum(iif(yjym='',1,0)),
sum(iif(yjym<>'',1,0))*100.0/(select count(jym)from xmjh where fa not in ('1-下架交易','5-移出柜面系统'))
from xjdz
where tcrq<=date('now')
