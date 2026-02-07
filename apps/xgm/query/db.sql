-- bug 清单
CREATE TABLE if not exists wtqd(
    bh      text    primary key, -- Bug编号
    bt      text,   --  Bug标题
    yzcd    text,   --	严重程度
    cxbz    text,   --	重现步骤
    zt      text,   --	Bug状态
    jhcs    int,    --  激活次数
    cjr     text,   --	由谁创建
    cjrq    text,   --	创建日期
    zp      text,   --	指派给
    jjr     text,   --	解决者
    jjfa    text,   --	解决方案
    jjrq    text,   --	解决日期
    gbr     text,   --  由谁关闭
    gbrq    text    --	关闭日期
);
-- 项目计划表
CREATE TABLE if not exists xmjh(
    jym     text    primary key,    -- 交易码
    jymc    text,   --  交易名称
    jyz     text,   -- 交易组
    jyzm    text,   -- 交易组名
    yjcd    text,   -- 一级菜单
    ejcd    text,   -- 二级菜单
    bs      int,    -- 交易笔数
    lx      text,   -- 类型：0-本部门，1-总行部门，2-分行特色
    ywbm    text,   -- 业务部门
    zx      text,   -- 中心
    lxr     text,   -- 业务联系人
    fa      text,   -- 改造方案  0-下架交易，1-直接迁移，2-改造迁移，3-重新设计,4-移出柜面系统
    pc      text,   -- 批次
    sfwc    text,   -- 是否完成
    bz      text,   -- 备注信息
    xjym    text    -- 新交易码
);

CREATE INDEX if not exists xmjh_ywbm on xmjh(ywbm);
CREATE INDEX if not exists xmjh_zx on xmjh(zx);
CREATE INDEX if not exists xmjh_lxr on xmjh(lxr);
CREATE INDEX if not exists xmjh_sfwc on xmjh(sfwc);
CREATE INDEX if not exists xmjh_fa on xmjh(fa);

-- 需求明细表
CREATE TABLE if not exists xqmxb(
    jym     text,   --交易码
    jymc    text,   --交易名称
    bxbm    text,   --提交部门
    bxr     text,   --提交人
    xqmc    text,   --需求名称
    zt      text,   --状态
    jhyf    text,   --应提交月份
    tjrq    text,   --实际提交
    psrq    text,   --需求评审
    zstjrq  text,   --提交开发日期
    bz      text    --备注
);
CREATE INDEX if not exists xqmxb_jym on xqmxb(jym);
CREATE INDEX if not exists xqmxb_bxbm on xqmxb(bxbm);
CREATE INDEX if not exists xqmxb_tjrq on xqmxb(tjrq);
CREATE INDEX if not exists xqmxb_psrq on xqmxb(psrq);
CREATE INDEX if not exists xqmxb_zstjrq on xqmxb(zstjrq);

-- 问题故障表
CREATE TABLE if not exists wtgzb(
    xh             text,       -- 序号
    jygn           text,       -- 交易/功能
    wtms           text,       -- 问题描述
    yzx            text,       -- 严重性
    tcrq           text,       -- 提出日期
    tcjg           text,       -- 提出机构
    yhfx           text,       -- 原因分析
    zt             text,       -- 状态
    wtfl           text,       -- 问题分类
    clfa           text,       -- 处理方案
    jhbb           text,       -- 计划版本
    zrr            text,       -- 责任人
    bz             text       -- 备注
);
-- 分支行业务专家
CREATE TABLE if not exists ywzj(
    fh             text,       -- 分行
    xm             text,       -- 姓名
    gh             text,       -- 工号
    kh             text,       -- 卡号
    xb             text,       -- 性别
    cdsj           text,       -- 抽调时间
    bdsj           text,       -- 报到时间
    lxr            text,       -- 总行联系人
    gznr           text       -- 工作内容
);
-- 故障中断表
CREATE TABLE if not exists bkjl(
    rq             text,       -- 日期
    jgh            text,       -- 机构号
    fh             text,       -- 所属分行
    ip             text,       -- IP地址
    czxt           text,       -- 操作系统
    nc             text,       -- 内存
    zdcs           text,       -- AB3中断次数
    primary key(rq,ip)
);
-- 开发计划表
CREATE TABLE if not exists kfjh(
    jym            text,       -- 交易码
    xqzt           text,       -- 需求状态
    kfzt           text,       -- 开发状态
    jhbb           text,       -- 计划版本
    jsxz           text,       -- 开发组
    kjfzr          text,       -- 行方负责人
    kfzz           text,       -- 技术组长
    qdkf           text,       -- 前端开发
    hdkf           text,       -- 后端开发
    lckf           text,       -- 流程开发
    kskf           text,       -- 开始开发
    wckf           text,       -- 完成开发
    wccs           text,       -- 完成集成测试
    wcys           text,       -- 完成验收
    primary key(jym)
);
-- 新旧对照表
CREATE TABLE if not exists xjdz(
    jym            text,       -- 交易码
    jymc           text,       -- 交易码称
    yjym           text,       -- 原交易码
    yjymc          text,       -- 原交易码称
    tcrq           text,       -- 投产日期
    zs             text,       -- 状态
    bz             text       -- 备注
);
CREATE TABLE if not exists loadfile(
    name	text,
    path	text,
    mtime	text,
    ver	text,primary key(name,path)
);

-- 分工明细表
create table if not exists fgmxb(
    bh  text   primary key,--验收条目编号
    ywxz    text,-- 业务小组
    jsxz    text -- 技术小组
);

-- 交易对照表
create table if not exists jydzb(
    bh  text,--验收条目编号
    jym  text,   -- 新交易码
    yjym text,   -- 交易码
    primary key(bh,yjym)
);

-- 验收条目表
create table if not exists ystmb(
    tcrq    text,   -- 投产日期
    bh      text primary key, -- 验收条目编号
    mc      text,   -- 功能名称
    gs      text,   -- 功能概述
    glxt    text,   --关联系统
    jsfzr  text,   -- 技术负责人
    zt      text,   -- 状态
    wbce    text,   -- 外包测试
    hxzc    text,   -- 核心支持人员
    cszt    text,   -- 测试状态：0-冒烟不通过（代码未同步、漏资源、无交易码、开发中、主流程不通）；1-冒烟测试中；2-集成测试中；3-待验收；4-验收完成
    bz      text,   -- 备注
    ywry    text   -- 业务人员
);

-- 人员清单
create table if not exists xmryb(
    xz  text,-- 小组
    lb text,    -- 类别：0-总行科技，1-通用外包，2-赞同公司，3-分行业务
    xm text     -- 姓名
);
create index if not exists xmryb_xz on xmryb(xz);
create index if not exists xmryb_lb on xmryb(lb);
create index if not exists xmryb_xm on xmryb(xm);

create table if not exists bbap(
    jhbb    text,
    tcrq    text,
    wcys    text
);

-- 缺陷指标
create table if not exists qxzb(
    qxbh   text primary key,  -- 缺陷编号
    xmbh    text,   -- 项目编号
    tcbb    text,   -- 投产版本
    xmjlgh  text,   -- 项目经理工号
    xmjl    text,   -- 项目经理
    csjlgh  text,   -- 测试经理工号
    csjl    text,   -- 测试经理
    csbm    text,   -- 测试部门
    yzcd    text,   -- 严重程度
    zt      text,   -- bug状态
    sfyq    text,   -- 是否延期
    xt      text,   -- 系统
    jhcs    int,    -- 激活次数
    cjrgh   text,   -- 创建人工号
    cjr     text,   -- 创建人
    jjrgh   text,   -- 解决人工号
    jjr     text,   -- 解决人
    csjlczsc  real,   -- 测试经理操作时长
    xmjlczsc  real,   -- 项目经理操作时长
    kfryczsc  real,   -- 开发人员操作时长
    csryczsc  real,   -- 测试人员操作时长
    yzbtgcs     int,    -- 验证不通过次数
    cjsj    text   -- 创建时间
);

-- 下架交易清单
create table if not exists xjjy(
    jym         text primary key,   --交易码
    xjfs        text,   -- 下架方式
    xjrq        text    -- 下架日期
);
