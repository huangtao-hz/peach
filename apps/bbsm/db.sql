 CREATE TABLE if not exists bbsm(
    rq      text,   -- 日期
    xm      text,   -- 系统或项目
    jym     text,   -- 交易码
    jymc    text,   -- 交易名称
    nr      text,   -- 测试内容
    yhyy    text,   -- 优化原因
    yzjg    text,   -- 验证机构
    wcsj    text,   -- 完成时间
    yzyq    text,   -- 验证要求
    lxr     text,   -- 联系人
    yzsj    time    -- 验证时间
);
CREATE INDEX if not exists bbsm_jym on bbsm(jym);
CREATE INDEX if not exists bbsm_rq on bbsm(rq);
