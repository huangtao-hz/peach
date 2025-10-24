-- 开发计划统计表，列表显示开发计划
select a.jhbb,count(a.jym),count(a.jym)*100.0/(select count(jym)from xmjh where fa not in ("1-下架交易","5-移出柜面系统"))
from kfjh a
left join xmjh b
on a.jym=b.jym
where b.sfwc not like "5%"
group by jhbb
union
select "合计",(select count(jym)
from xmjh
where fa not in ("1-下架交易","5-移出柜面系统") and sfwc not like "5%"),
(select count(jym)from xmjh
where fa not in ("1-下架交易","5-移出柜面系统") and sfwc not like "5%")*100.0/(select count(jym)
from xmjh where fa not in ("1-下架交易","5-移出柜面系统"))
order by jhbb
