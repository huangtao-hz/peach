-- 根据验收明细表更新开发状态
update kfjh
set kfzt=ystmb.cszt
from jydzb,ystmb
where kfjh.jym=jydzb.yjym and ystmb.bh=jydzb.bh and ystmb.bh is not null and kfzt<>ystmb.cszt
