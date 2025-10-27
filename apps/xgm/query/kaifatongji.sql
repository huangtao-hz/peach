-- 开发统计，统计各开发小组当期的完成情况
select b.jsxz,sum(iif(a.cszt like '0%',1,0)),
sum(iif(a.cszt like '1%',1,0)),
sum(iif(a.cszt like '2%',1,0)),
sum(iif(a.cszt like '3%',1,0)),
sum(iif(a.cszt like '4%',1,0)),
count(a.bh) as sl
from ystmb a left join fgmxb b on a.bh=b.bh
where tcrq=?
group by b.jsxz
order by sl desc
