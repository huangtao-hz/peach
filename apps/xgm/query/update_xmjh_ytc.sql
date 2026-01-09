update xmjh
set sfwc='5-已投产'
from xjdz
where sfwc<>'5-已投产' and xmjh.jym=xjdz.yjym and tcrq<date('now','localtime') and tcrq<>''
