-- 根据计划版本更新开发计划的时间
update kfjh
set
kskf = date(kskf,printf("%d days",julianday(bbap.wcys)-julianday(kfjh.wcys))),
wckf = date(wckf,printf("%d days",julianday(bbap.wcys)-julianday(kfjh.wcys))),
wccs = date(wccs,printf("%d days",julianday(bbap.wcys)-julianday(kfjh.wcys))),
wcys = bbap.wcys
from bbap
where kfjh.jhbb=bbap.jhbb and kfjh.wcys<>bbap.wcys
