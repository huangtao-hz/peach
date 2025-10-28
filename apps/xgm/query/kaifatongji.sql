-- 开发统计，统计各开发小组当期的完成情况
select ifnull(b.xz,"其他") as jsxz,
c.rs,
sum(iif(a.cszt like '0%',1,0)),
sum(iif(a.cszt like '1%',1,0)),
sum(iif(a.cszt like '2%',1,0)),
sum(iif(a.cszt like '3%',1,0)),
sum(iif(a.cszt like '4%',1,0)),
count(a.bh) as sl,
sum(d.xjys),sum(d.ljys)
from ystmb a
left join xmryb b on a.jsfzr=b.xm and lb="0-总行科技"
left join (select xz,count(xz)as rs from xmryb where lb in("0-总行科技","1-通用外包","2-赞同公司")group by xz)c on b.xz=c.xz
left join (select bh,count(distinct jym)as xjys,count(distinct yjym)as ljys from jydzb group by bh) d on a.bh=d.bh
where tcrq=?
group by jsxz
order by sl desc
