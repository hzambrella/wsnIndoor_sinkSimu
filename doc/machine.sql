-- 机器码分销项目

-- 所有金额的单位为厘(小数点后四位)

Drop table if exists machine_user_info;
-- 机器码分销用户表
CREATE TABLE machine_user_info (
    -- 用户id
    userid BIGINT(11) NOT NULL AUTO_INCREMENT,
    -- 用户角色
    role_code VARCHAR(32) NOT NULL,
    -- 用户父节点，0代表根。
    parentid BIGINT(20) NOT NULL DEFAULT '0',
    -- 用户账号
    account VARCHAR(25)NOT NULL,
    -- 密码
    pass VARCHAR(128) NOT NULL,
    -- 手机号，账号就是手机号，但是防止变更预留这个字段
    mobile VARCHAR(32) NOT NULL DEFAULT '',
     -- 真实名字
    real_name VARCHAR(32)NOT NULL DEFAULT '',
    -- 身份证号
    id_card VARCHAR(32) NOT NULL DEFAULT '',
    -- 预设置的银行卡号
    bank_card VARCHAR(32) NOT NULL DEFAULT '',
    -- 预设的地址信息
    address  VARCHAR(128) NOT NULL DEFAULT '',
    -- 状态
    -- 0:正常  1:冻结或删除  2：审核中
    status VARCHAR(45) NOT NULL DEFAULT '0',
    -- 创建时间
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 更新时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE(account),
    PRIMARY KEY (userid)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 AUTO_INCREMENT=100000; 
-- PARTITION BY HASH(userid DIV 1000000) PARTITIONS 10;
CREATE INDEX machine_user_info_idx0 ON machine_user_info(role_code);
CREATE INDEX machine_user_info_idx1 ON machine_user_info(parentid);

-- 管理员
INSERT INTO machine_user_info (userid,pass, mobile, account, real_name, role_code, parentid,id_card) VALUES (100000,'$2a$10$Ul7QTKp2VGtqmGME0Kj/ruMecFaODFtt7ONZUgaWvGJw1OuHui6Gy', '13907390836', '13907390836', '管理员', 'ADMIN', '0','0');

Drop table if exists machine_role;
-- 用户角色表
CREATE TABLE machine_role (
    -- 角色信息码。系统管理员填写这个code,尽量简单明了，一目了然。
    role_code VARCHAR(8) NOT NULL,
    -- 名字
    role_name VARCHAR(48) NOT NULL,
    -- 描述
    detail VARCHAR(256) NOT NULL  DEFAULT '',
    -- 父身份
    parent_role VARCHAR(8)NOT NULL DEFAULT 'ROOT',
    -- 创建时间
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 更新时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 状态
    -- 0:正常  1:冻结或删除 
    status INT(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (role_code)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8;
INSERT INTO machine_role (role_code, role_name, detail) VALUES ('ADMIN','平台管理员','系统管理员，可以查看所有角色，无需审核地创建大多数角色，查看所有信息');
INSERT INTO machine_role (role_code, role_name, detail,parent_role) VALUES ('CB','消费商','消费商','CKZX');
INSERT INTO machine_role (role_code, role_name, detail,parent_role) VALUES ('CITY','城市运营中心','城市运营中心','ADMIN');
INSERT INTO machine_role (role_code, role_name, detail,parent_role) VALUES ('MANAGER','管理中心','管理中心','CITY');
INSERT INTO machine_role (role_code, role_name, detail,parent_role) VALUES ('SERVER','服务运营中心','服务运营中心','MANAGER');
INSERT INTO machine_role (role_code, role_name, detail,parent_role) VALUES ('CKZX','创客中心','创客中心','SERVER');
INSERT INTO machine_role (role_code, role_name, detail,parent_role) VALUES ('CONSUMER','消费者','消费者','CKZX');
INSERT INTO machine_role (role_code, role_name, detail,parent_role) VALUES ('ROOT','系统管理员','系统管理员，用来创建角色，行为，维护系统','');

/*
Drop table if exists machine_action;
-- 行为表
CREATE TABLE machine_action (
    -- 行为码。系统管理员填写这个code,尽量简单明了，一目了然。 
    action_code VARCHAR(16) NOT NULL,
    -- 名字
    acrtion_name VARCHAR(48) NOT NULL,
    -- 描述
    detail VARCHAR(256) NOT NULL DEFAULT '',
    -- 创建时间
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 更新时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 状态
    -- 0:正常  1:冻结或删除 
    status INT(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (action_code)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8;
INSERT INTO machine_action (action_code, acrtion_name, detail) VALUES ('ROOT', 'root', 'root行为');
INSERT INTO machine_action (action_code, acrtion_name, detail) VALUES ('ADMIN', '平台管理', '默认的平台管理行为');
INSERT INTO machine_action (action_code, acrtion_name, detail) VALUES ('S_CONSUME', '查看消费者', '查看消费者行为');
INSERT INTO machine_action (action_code, acrtion_name, detail) VALUES ('C_CONSUME', '创建消费者', '创建消费者行为');
INSERT INTO machine_action (action_code, acrtion_name, detail) VALUES ('U_CONSUME', '更改消费者', '更改消费者行为');
INSERT INTO machine_action (action_code, acrtion_name, detail) VALUES ('D_CONSUME', '删除消费者', '删除消费者行为');
*/

Drop table if exists machine_role_act;
-- 角色-行为基本表
-- 采用json编码的string数组存储权限，数组内容是role_code
-- 如果数组有ALL，则除了数组中的role_code，其它身份都有权限进行操作
-- 基本表包含了增删查改操作的权限。其它行为可以建立拓展表
CREATE TABLE machine_role_act (
    -- 角色码
    role_code VARCHAR(32) NOT NULL,
    -- 创建时间
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 更新时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 状态
    -- 0:正常  1:冻结或删除
    status INT(11) NOT NULL DEFAULT '0',
    -- create_role
    create_role VARCHAR(256) NOT NULL DEFAULT '[]', 
    -- update_role
    update_role VARCHAR(256) NOT NULL DEFAULT '[]', 
    -- delete_role
    delete_role VARCHAR(256) NOT NULL DEFAULT '[]', 
    -- select_role
    select_role VARCHAR(256) NOT NULL DEFAULT '[]', 
    -- 可分配机器码给的角色
    distribute_role VARCHAR(256) NOT NULL DEFAULT '[]', 
    -- 是否可销售实物。若是，会给消费者生成链接
    sale_role BOOLEAN NOT NULL DEFAULT '0',
    -- 是否可下订单购买实物。若是，则会进入到下单页面
    purchase_role BOOLEAN NOT NULL DEFAULT '0',
    PRIMARY KEY (role_code)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8;

-- 系统管理员，行为是系统默认的系统管理行为
INSERT INTO machine_role_act (role_code, create_role, update_role, delete_role, select_role,distribute_role)
 VALUES ('ROOT', '[\"ALL\"]', '[\"ALL\",\"CB\"]', '[\"ALL\"]', '[\"ALL\"]','[\"ADMIN\"]');

-- 平台管理员，行为是系统默认的平台管理行为
INSERT INTO machine_role_act (role_code, create_role, update_role, delete_role, select_role,distribute_role) 
VALUES ('ADMIN', '[\"ALL\",\"ROOT\",\"CB\"]','[\"ALL\",\"ROOT\"]', '[\"ALL\",\"ROOT\"]', '[\"ALL\",\"ROOT\"]','[\"ALL\",\"ROOT\"]');

-- 城市运营中心
INSERT INTO machine_role_act (role_code, create_role, update_role, delete_role, select_role,distribute_role) 
VALUES ('CITY', '[\"MANAGER\",\"SERVER\",\"CKZX\"]', '[\"MANAGER\",\"SERVER\",\"CKZX\"]',  '[\"MANAGER\",\"SERVER\",\"CKZX\"]', '[\"MANAGER\",\"SERVER\",\"CKZX\"]', '[\"MANAGER\",\"SERVER\",\"CKZX\"]');

-- 管理中心
INSERT INTO machine_role_act (role_code, create_role, update_role, delete_role, select_role,distribute_role) 
VALUES ('MANAGER', '[\"SERVER\",\"CKZX\"]', '[\"SERVER\",\"CKZX\"]',  '[\"SERVER\",\"CKZX\"]', '[\"SERVER\",\"CKZX\"]', '[\"SERVER\",\"CKZX\"]');

-- 服务运营中心
INSERT INTO machine_role_act (role_code, create_role, update_role, delete_role, select_role,distribute_role) 
VALUES ('SERVER', '[\"CKZX\"]','[\"CKZX\"]', '[\"CKZX\"]', '[\"CKZX\"]','[\"CKZX\"]');

-- 创客中心，查询和删除role_code为CONSUME的行为
INSERT INTO machine_role_act (role_code, create_role, update_role, delete_role, select_role,distribute_role,sale_role) 
VALUES ('CKZX', '[\"CONSUMER\",\"CKZX\"]', '[\"CONSUMER\",\"CKZX\",\"CB\"]', '[\"CONSUMER\",\"CKZX\",\"CB\"]', '[\"CONSUMER\",\"CKZX\",\"CB\"]', '[\"CONSUMER\",\"CKZX\",\"CB\"]','1');

-- 消费者
INSERT INTO machine_role_act (role_code, create_role, update_role, delete_role, select_role,distribute_role,sale_role,purchase_role) 
VALUES ('CONSUMER', '[\"CONSUMER\"]', '[\"CONSUMER\"]', '[]', '[\"CONSUMER\"]','[\"CONSUMER\"]',1,1);

-- 消费商
INSERT INTO machine_role_act (role_code, create_role, update_role, delete_role, select_role,distribute_role,sale_role,purchase_role) 
VALUES ('CB', '[\"CONSUMER\"]', '[\"CONSUMER\",\"CB\"]', '[]', '[\"CONSUMER\",\"CB\"]','[\"CONSUMER\",\"CB\"]',1,1);

Drop table if exists machine_bank_info;
-- 用户银行卡信息
CREATE TABLE machine_bank_info (
    -- 用户ID
    userid BIGINT(20) NOT NULL,
    -- 银行卡号
    card_sn VARCHAR(64) NOT NULL,
    -- 绑定卡的身份证姓名
    idcard_name VARCHAR(64) NOT NULL,
    -- 绑定卡的身份证号码
    idcard_num VARCHAR(18) NOT NULL,
    -- 开户银行
    head_bank VARCHAR(64) NOT NULL,
    -- 开户支行
    sub_bank VARCHAR(128) NOT NULL DEFAULT '',
    -- 银行预留的手机号码
    mobile VARCHAR(32) NOT NULL,
    -- 创建时间
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 更新时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 状态
    -- 0:正常  1:冻结或删除
    status INT(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (userid)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8;

Drop table if exists machine_addr_info;
-- 用户地址
CREATE TABLE machine_addr_info (
    -- 地址id
    addr_id INT(11) NOT NULL,
    -- 用户id
    userid BIGINT(20) NOT NULL,
    -- 用户地址
    addr VARCHAR(128) DEFAULT NULL,
    -- 创建时间
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 更新时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 状态
    -- 0:正常  1:冻结或删除
    status INT(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (addr_id)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8;

Drop table if exists machine_promoter_money;
-- 用户佣金表，当用户审核通过后创键，状态随用户状态而改变。单位:元后小数点4位
CREATE TABLE machine_promoter_money (
    -- 用户ID
    userid BIGINT(20) NOT NULL AUTO_INCREMENT,
    -- 可提现佣金
    cash INT(11) NOT NULL DEFAULT 0,
    -- 拥金余额(压款+可提现金额，)
    -- balance INT NOT NULL DEFAULT '0',
    -- 申请提现中
    apply INT NOT NULL DEFAULT 0,
    -- 已提现佣金
    withdraw INT(11) NOT NULL DEFAULT 0,
    -- 佣金总额=可提现佣金+申请提现中+已提现佣金
    total INT(11) DEFAULT 0,
     -- 创建时间
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      -- 更新时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (userid)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8;

INSERT INTO machine_promoter_money (userid,cash,total) VALUES ('100000','4000000','4000000');


Drop table if exists machine_money_record;
-- 佣金流水表（获取和提现）结算和提现的记录都在此。
CREATE TABLE machine_money_record (
    -- 订单id
    order_id VARCHAR(64) NOT NULL,
    -- 用户ID
    userid BIGINT(20) NOT NULL,
    -- operateId 经办人id,提现时审核的人记录于此。
    operateid BIGINT(20) NOT NULL,
    -- 金额(单位:元后小数点4位), 正数为提现
    amount INT(11) NOT NULL DEFAULT '0',
    -- 操作 APPLY代表申请提现，WITHDREW代表已提现，IN代表佣金进账。
    operate VARCHAR(10) NOT NULL DEFAULT 'OUT',
    -- 
    -- 提现渠道，0微信号, 1银行卡提现
    -- pay_channel INT(11) DEFAULT NULL,
    -- 提现，0时，此值为微信的appid; 1时为银行类别
    -- appid_item VARCHAR(256) NOT NULL DEFAULT '',
    -- 提现，0时，此值为微信的appid; 1时为银行类别
    -- openid_card VARCHAR(256) NOT NULL DEFAULT '',
    -- 个税百分之几, 只存整数部分, 如:20%存储为20
    -- tax INT(11) NOT NULL DEFAULT '0',
    -- 税款
    -- cash_tax BIGINT(20) NOT NULL DEFAULT '0',
    -- 实收现金
    -- cash_pay BIGINT(20) NOT NULL DEFAULT '0',
    -- 创建时间
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 更新时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 订单状态
    -- 0,审核中; 1，已到账，2，已取消, 填写备注值
    status INT(11) NOT NULL DEFAULT '0',
    -- 备注。获取时，将机器订单号填写于此。
    memo VARCHAR(256)NOT NULL DEFAULT '',
    PRIMARY KEY (order_id),
    KEY machine_promoter_withdraw_idx0 (userid),
    KEY machine_promoter_withdraw_idx1 (create_time)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8;
-- PARTITION BY KEY (order_id) PARTITIONS 10;

Drop table if exists machine_distribute_record;
-- 机器码分配记录表，兼机器码归属表
CREATE TABLE machine_distribute_record (
    -- 机器码
    machine_code VARCHAR(45) NOT NULL,
    -- 发放者 （0代表根,机器码的最初拥有者）
    from_userid BIGINT(20) NOT NULL DEFAULT '0',
    -- 分配者（机器码的拥有者） 
    to_userid BIGINT(20) NOT NULL DEFAULT '100000',
    -- 是否是机器码的当前拥有者
    -- 1是,0否 转化为
    is_owner BOOLEAN NOT NULL DEFAULT '1',
    -- 0：正常 1：售出状态，机器码售出具体状态可查看订单 
    -- 2:待结算，表示这条记录“分配者”佣金正在结算
    -- 3：已结算，表示这条记录的“分配者”的佣金已经结算
    -- 4：被删除,冻结，失效
    status INT(11) NOT NULL DEFAULT '0',
     -- 创建时间，若需机器码查询走向，可以用这个排序
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (from_userid,to_userid,machine_code)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8;
CREATE INDEX machine_distribute_record_idx0 ON machine_distribute_record(from_userid);
CREATE INDEX machine_distribute_record_idx1 ON machine_distribute_record(to_userid);
CREATE INDEX machine_distribute_record_idx2 ON machine_distribute_record(machine_code);

Drop table if exists  machine_sale_record;
-- 机器订单和交易记录表
CREATE TABLE machine_sale_record (
    order_id VARCHAR(64) NOT NULL,
    -- 卖家
    seller_id BIGINT(20) NOT NULL,
    -- 买家
    purchase_id BIGINT(20) NOT NULL,
    -- 机器码
    machine_code VARCHAR(45) NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 0：已下单 1：已发件 2：已到货
    -- 3：已收货 4：已付款 5：退货中 6：已退货 7:被删除，冻结
    status INT(11) NOT NULL DEFAULT '0',
    -- 地址信息
    addr VARCHAR(128) NOT NULL,
    -- 电话信息
    mobile VARCHAR(32) NOT NULL,
    -- 银行卡信息
    bank_card VARCHAR(32) NOT NULL DEFAULT '',
    -- 备注
    memo VARCHAR(256) NOT NULL DEFAULT '',
    PRIMARY KEY (order_id,machine_code)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 PARTITION BY KEY (order_id) PARTITIONS 10;
CREATE INDEX machine_sale_record_idx0 ON machine_sale_record(create_time);
CREATE INDEX machine_sale_record_idx1 ON machine_sale_record(seller_id);
CREATE INDEX machine_sale_record_idx2 ON machine_sale_record(purchase_id);


-- 查询所有子节点
-- SELECT * FROM machine_user_info 
-- WHERE FIND_IN_SET(userid,queryAllChild(?))and userid!=? and role_code=?
DROP function IF EXISTS queryAllChild;

DELIMITER $$

CREATE FUNCTION queryAllChild(areaId INT)
RETURNS VARCHAR(4000)
BEGIN
DECLARE sTemp VARCHAR(4000);
DECLARE sTempChd VARCHAR(4000);

SET sTemp='$';
SET sTempChd = CAST(areaId AS CHAR);

WHILE sTempChd IS NOT NULL DO

SET sTemp= CONCAT(sTemp,',',sTempChd);
SELECT 
    GROUP_CONCAT(userid)
INTO sTempChd FROM
    machine_user_info
WHERE
    FIND_IN_SET(parentid, sTempChd) > 0 ;


END WHILE;
RETURN sTemp;
END$$

DELIMITER ;

-- 寻找父节点
-- 例：SELECT * FROM hz_test.machine_user_info 
-- where FIND_IN_SET(userid,queryChildrenAreaInfo1(4))
-- and userid!=4 and role_code="ADMIN;
DROP FUNCTION IF EXISTS queryAllParent;
DELIMITER $$
CREATE FUNCTION queryAllParent(areaId INT)
RETURNS VARCHAR(4000)
BEGIN
DECLARE sTemp VARCHAR(4000);
DECLARE sTempChd VARCHAR(4000);

SET sTemp='$';
SET sTempChd = CAST(areaId AS CHAR);
SET sTemp = CONCAT(sTemp,',',sTempChd);

SELECT parentid INTO sTempChd FROM machine_user_info WHERE userid = sTempChd;
WHILE sTempChd <> 0 DO
SET sTemp = CONCAT(sTemp,',',sTempChd);
SELECT parentid INTO sTempChd FROM machine_user_info WHERE userid = sTempChd;
END WHILE;
RETURN sTemp;
END$$

DELIMITER ;



-- 机器码流向查询
-- 例：
-- SELECT * FROM machine_distribute_record WHERE FIND_IN_SET(to_userid,queryMachineFlow(2)) 
-- and machine_code='abc-1' order by create_time desc
DROP function IF EXISTS queryMachineFlow;

DELIMITER $$

CREATE FUNCTION queryMachineFlow(areaId INT)
RETURNS VARCHAR(4000)
BEGIN
DECLARE sTemp VARCHAR(4000);
DECLARE sTempChd VARCHAR(4000);

SET sTemp='$';
SET sTempChd = CAST(areaId AS CHAR);

WHILE sTempChd IS NOT NULL DO

SET sTemp= CONCAT(sTemp,',',sTempChd);
SELECT 
    GROUP_CONCAT(to_userid)
INTO sTempChd FROM
    machine_distribute_record
WHERE
    FIND_IN_SET(from_userid, sTempChd) > 0 ;


END WHILE;
RETURN sTemp;
END$$

DELIMITER ;

