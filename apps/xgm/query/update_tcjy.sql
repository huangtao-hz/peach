-- 根据项目计划更新投产交易的联系人
update xjdz set
ywbm = xmjh.ywbm,
lxr = xmjh.lxr
from xmjh
where xjdz.ywbm= "" and xjdz.yjym <> "" and xjdz.yjym = xmjh.jym
