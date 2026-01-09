update xmjh
set sfwc="3-已完成开发"
from ystmb,jydzb
where xmjh.jym=jydzb.yjym
and jydzb.bh=ystmb.bh
and ystmb.cszt in ("2-集成测试中","3-待验收","4-验收完成")
and sfwc not in ("3-已完成开发","5-已投产")
