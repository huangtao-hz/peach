-- 业务统计，统计各业务小组当期验收条目的完成情况
select b.ywxz,c.rs,
sum(iif(a.cszt like '0%',1,0)),
sum(iif(a.cszt like '1%',1,0)),
sum(iif(a.cszt like '2%',1,0)),
sum(iif(a.cszt like '3%',1,0)),
sum(iif(a.cszt like '4%',1,0)),
count(a.bh) as sl,sum(d.xjys),sum(d.ljys)
from ystmb a left join fgmxb b on a.bh=b.bh
left join (select xz,count(xz)as rs from xmryb where lb="3-分行业务" group by xz)c on b.ywxz=c.xz
left join (select bh,count(distinct jym)as xjys,count(distinct yjym)as ljys from jydzb group by bh) d on a.bh=d.bh

where tcrq=?
group by b.ywxz
order by sl desc
