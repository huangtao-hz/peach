update xmjh
set xjym=printf("%s-%s",xjdz.jym,xjdz.jymc)
from xjdz
where xmjh.jym=xjdz.yjym and xjym=""
