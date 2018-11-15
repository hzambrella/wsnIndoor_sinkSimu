package machine

// 格式：每五个字段换行。
//当被连表数据不存在时，采用left join才能出现结果
const (
	//检查账号
	checkUserByNameSql = `
SELECT
    status
FROM
	machine_user_info
WHERE
	account=?
`
	//登录检查
	checkLoginSql = `
SELECT
    userid,pass,status
FROM
	machine_user_info
WHERE
	account=?
	`
	//获取用户基本信息，通过 ID
	getUserBaseInfoByIdSql = `
SELECT
	tb1.userid,tb1.role_code,tb1.parentid,tb1.mobile,
	tb1.account,tb1.real_name,tb1.id_card,tb1.bank_card,tb1.address,
	tb1.status,tb1.create_time,tb1.update_time,
    ifnull((SELECT account FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_name,
	ifnull((SELECT real_name FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_real_name,
	tb2.role_name as role_name
FROM
	machine_user_info as tb1
JOIN 
	machine_role as tb2
ON 
	tb1.role_code=tb2.role_code
WHERE
	userid=?
	`

	//获取佣金信息，通过ID
	getUserBalanceInfoByIdSql = `
SELECT
	 cash,apply,withdraw,total,
	 create_time,update_time
FROM
	machine_promoter_money
WHERE
	userid=?
	`

	//通过账号获得基本信息
	getUserBaseInfoByAccountSql = `
SELECT
	tb1.userid,tb1.role_code,tb1.parentid,tb1.mobile,
	tb1.account,tb1.real_name,tb1.id_card,tb1.bank_card,tb1.address,
	tb1.status,tb1.create_time,tb1.update_time,
    ifnull((SELECT account FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_name,
	ifnull((SELECT real_name FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_real_name,
	tb2.role_name as role_name
FROM
	machine_user_info as tb1
JOIN 
	machine_role as tb2
ON 
	tb1.role_code=tb2.role_code
WHERE
	tb1.account=?
	`

	//通过账号获得佣金信息
	getUserBalanceInfoByAccountSql = `
SELECT
	 tb1.cash,tb1.apply,tb1.withdraw,tb1.total,
	 tb1.create_time,tb1.update_time
FROM
	machine_promoter_money as tb1
JOIN
	machine_user_info as tb2
ON
	tb1.userid=tb2.userid
WHERE
	tb2.account=?
	`

	//创建基本信息
	createUserBaseInfoSql = `
INSERT INTO 
	machine_user_info 
	(parentid,pass, mobile, account, real_name, 
		role_code ,id_card,bank_card,status) 
VALUES 
	(?,?,?,?,?,
	?,?,?,?);
	`

	//创建佣金信息
	createUserBalanceInfoSql = `
INSERT INTO 
	machine_promoter_money 
	(userid) 
VALUES 
	(?);
	`
	//获得直接父节点信息
	getParentNodeInfoSql = `
SELECT
	 userid,role_code,parentid,mobile,
	 account,real_name,id_card,bank_card,address,
	 status,create_time,update_time
FROM
	machine_user_info
WHERE
	userid=(SELECT parentid FROM machine_user_info WHERE userid=?)
	`
	checkChildNodeInfoSql = `
SELECT
	count(*)
FROM
	machine_user_info
WHERE
	parentid=? and 	userid=?
	`

	getChildNodeInfoSql = `
SELECT
	tb1.userid,tb1.role_code,tb1.parentid,tb1.mobile,tb1.account,
	tb1.real_name,tb1.id_card,tb1.bank_card,tb1.address,tb1.status,
	tb1.create_time,tb1.update_time,
	ifnull((SELECT tb1.account FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_name,
	ifnull((SELECT tb1.real_name FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_real_name,
	ifnull(tb3.cash,0) as cash,
    ifnull(tb3.apply,0) as apply,
    ifnull(tb3.withdraw,0)as withdraw,
    ifnull(tb3.total,0) as total
FROM
	machine_user_info as tb1
LEFT JOIN
	machine_promoter_money as tb3
ON
	tb1.userid=tb3.userid
WHERE
	tb1.parentid=?
`
	//更改用户基本信息
	updateUserBaseInfoSql = `
UPDATE
	machine_user_info
SET
	mobile=?,account=?,real_name=?,id_card=?,bank_card=?,
	address=?
WHERE
	userid=?
	`

	//更改用户角色信息
	updateUserRoleSql = `
UPDATE
	machine_user_info
SET
	role_code=?
WHERE
	userid=?
	`

	//更改用户佣金信息
	updateUserBalanceInfoSql = `
UPDATE
	machine_promoter_money 
SET
	cash=?,apply=?,withdraw=?,total=?
WHERE
	userid=?
	`

	//更改用户状态
	updateUserStatusSql = `
UPDATE
	machine_user_info
SET
	status=?
WHERE
	userid=?
	`

	//更改用户密码
	updateUserPassSql = `
UPDATE
	machine_user_info
SET
	pass=?
WHERE
	userid=?
	`
	updateUserPassByAccountSql = `
UPDATE
	machine_user_info
SET
	pass=?
WHERE
	account=?
	`

	//获得子节点
	getAllChildNodeInfoSql = `
SELECT
	tb1.userid,tb1.role_code,tb1.parentid,tb1.mobile,tb1.account,
	tb1.real_name,tb1.id_card,tb1.bank_card,tb1.address,tb1.status,
	tb1.create_time,tb1.update_time,
	ifnull((SELECT tb1.account FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_name,
	ifnull((SELECT tb1.real_name FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_real_name,
	ifnull(tb3.cash,0) as cash,
    ifnull(tb3.apply,0) as apply,
    ifnull(tb3.withdraw,0)as withdraw,
    ifnull(tb3.total,0) as total
FROM
	machine_user_info as tb1
LEFT JOIN
	machine_promoter_money as tb3
ON
	tb1.userid=tb3.userid
WHERE
	FIND_IN_SET(tb1.userid,queryAllChild(?))and tb1.userid!=?
	`

	//获得子节点数量
	getAllChildNodeNumSql = `
SELECT
	count(userid)
FROM
	machine_user_info as tb1
WHERE
	FIND_IN_SET(tb1.userid,queryAllChild(?))and tb1.userid!=?
	`

	//获得父节点
	getAllParentNodeInfoSql = `
SELECT
	tb1.userid,tb1.role_code,tb1.parentid,tb1.mobile,tb1.account,
	tb1.real_name,tb1.id_card,tb1.bank_card,tb1.address,tb1.status,
	tb1.create_time,tb1.update_time,
	ifnull((SELECT tb1.account FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_name,
	ifnull((SELECT tb1.real_name FROM machine_user_info WHERE userid=tb1.parentid),'无') as parent_real_name,
	ifnull(tb3.cash,0) as cash,
    ifnull(tb3.apply,0) as apply,
    ifnull(tb3.withdraw,0)as withdraw,
    ifnull(tb3.total,0) as total
FROM
	machine_user_info as tb1
LEFT JOIN
	machine_promoter_money as tb3
ON
	tb1.userid=tb3.userid
WHERE
	FIND_IN_SET(tb1.userid,queryAllParent(?))and tb1.userid!=?

	`

	//获得角色名字
	getRoleNameSql = `
SELECT
	role_code,role_name
FROM
	machine_role
WHERE
	role_code=? and status=0
	`
	getAllRoleNameSql = `
SELECT
	role_code,role_name
FROM
	machine_role
WHERE
	status=0
	
	`
	//获得角色的某项操作信息
	getRoleByIdSql = `
SELECT
	%s
FROM
	machine_role_act as tb1
JOIN
	machine_role as tb2
ON 
	tb1.role_code=tb2.role_code
JOIN
	machine_user_info as tb3
ON
	tb2.role_code=tb3.role_code
WHERE
	tb3.userid=? and tb1.status=0

	`

	//密码
	getUserBase = `
SELECT
	pass
FROM
	machine_user_info
WHERE
	userid=?
`
	//查询用户的佣金流水记录。where拼接
	getUserMoneyRecordSql = `
SELECT
	order_id,userid,operateid,amount,operate,
	update_time,create_time,status,memo,
	CASE status
		WHEN 0 then '审核中'
		WHEN 1 then '已到账'
		WHEN 2 then '已取消'
	ELSE  '状态未知' 
	END as status_name
FROM
	machine_money_record
	`
	//修改用户的佣金流水记录信息。
	updateUserMoneyRecordSql = `
UPDATE
	machine_money_record
SET
	operateid=?,status=?,memo=?
WHERE
	order_id=?
	`
	//增加佣金流水记录
	createMoneyRecordSql = `
INSERT INTO 
	machine_money_record 
	(order_id,userid,operateid,amount,operate,
	status,memo) 
VALUES 
	(?,?,?,?,?,
	?,?);
	`
	//查看父身份
	getParentRoleSql = `
SELECT
	parent_role
FROM
	machine_role
WHERE
	role_code=?
	`

	mimashengjiQuerySql = `
SELECT
	userid,pass
FROM
	machine_user_info
	`
)
