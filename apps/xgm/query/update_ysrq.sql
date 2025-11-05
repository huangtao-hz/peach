update kfjh
set kskf=date(kskf,printf("%d days",julianday(bbap.wcys)-julianday(wcys))),
wckf=date(wckf,printf("%d days",julianday(bbap.wcys)-julianday(wcys))),
wccs=date(wccs,printf("%d days",julianday(bbap.wcys)-julianday(wcys))),
wcys=bbap.wcys
from bbap
where kfjh.jhbb=bbap.jhbb
