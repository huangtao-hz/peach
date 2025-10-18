select zx,
sum(iif((sfwc is null or sfwc ='0-尚未开始' or sfwc='' ),1,0)),       -- 未开始
sum(iif(sfwc in('1-已编写初稿','2-已提交需求/确认需规'),1,0)),       -- 已完成需求
sum(iif(sfwc in('3-已完成开发','4-已完成验收测试'),1,0)),       -- 开发中
sum(iif(sfwc = '5-已投产' ,1,0)) ,       -- 已完成需求
count(jym) as zs        -- 总数
from xmjh
where ywbm='运营管理部' and fa not in('1-下架交易','5-移出柜面系统')
group by zx
order by zs desc
