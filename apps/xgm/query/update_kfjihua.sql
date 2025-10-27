-- select a.jym,b.yjym,a.kfzt,c.cszt
--from kfjh a
--left join jydzb b on a.jym=b.yjym
--left join ystmb c on b.bh=c.bh
--where c.bh is not null

update kfjh
set kfzt=ystmb.cszt
from jydzb,ystmb
where kfjh.jym=jydzb.yjym and ystmb.bh=jydzb.bh and ystmb.bh is not null
