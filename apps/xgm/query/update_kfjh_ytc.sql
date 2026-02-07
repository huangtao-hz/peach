-- 根据项目计划表更新开发状态
update kfjh
set kfzt="5-已投产"
from xmjh
where xmjh.jym=kfjh.jym and xmjh.sfwc="5-已投产"
and kfzt <>"5-已投产"
