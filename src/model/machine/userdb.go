package machine

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"strconv"
	"strings"
	"time"

	"tool/datastore"
	"tool/errors"
	"tool/money"

	"github.com/satori/go.uuid"
	//"github.com/qiniu/log"
	"tool/bcrypt"
)

//用户
var (
	// 封装错误，方便route层针对不同的错误进行不同的处理。如仅仅是查找结果为空和数据库系统错误的处理方式就不一样
	UserDataNotFound   = errors.New("user data not found")
	RoleDataNotFound   = errors.New("role data not found")
	UserWrongAuthParam = errors.New("wrong auth param")
	UserWrongPass      = errors.New("user pass is not correct")
	UserFrozen         = errors.New("user is frozen")
	UserChecking       = errors.New("user is checking")
	UserNotEnoughMoney = errors.New("user is not enough withdrew money")
	UserNotEnoughAuth  = errors.New("user have not enouth authority to operate")
	formatTime         = "2006-01-02 15:04:05"
)

type UserDB interface {
	Close()
	//==登录==
	// 检查用户状态
	CheckAccount(account string) (bool, error)
	// 登录校验,返回用户id
	LoginCheck(account string, pass string) (bool, int64, error)

	//==用户信息==
	// 根据用户id查询用户基本详情信息
	GetUserBaseInfoById(userid int64) (*UserBaseInfo, error)
	// 根据用户id查询用户佣金详情信息。若不存在记录，会创建一条新记录。
	GetUserBalanceInfoById(userid int64) (*UserBalanceInfo, error)
	// 根据用户账户获取信息基本详情信息
	GetUserBaseInfoByAccount(account string) (*UserBaseInfo, error)
	// 根据用户账户获取信息佣金详情信息
	GetUserBalanceInfoByAccount(account string) (*UserBalanceInfo, error)
	//查询直接上级信息
	GetParentNodeDirect(userid int64) (*UserBaseInfo, error)
	//查询上级信息,若role_code为空，就是上级的所有角色
	//status:用户状态。-1时是所有的。
	//acount:用户账户。若用户账户为空，就不筛选
	GetParentNodeInfo(userid int64, role_code string, status int) ([]UserInfo, error)
	//查询所有下级用户,若role_code为空，就是下级的所有角色
	//status:用户状态。-1时是所有的。
	//acount:用户账户。若用户账户为空，就不筛选
	//返回,数据总条数
	GetChildNodeNum(userid int64, role_code string, status int, account string, page, limit int)(int,error)
	//返回列表数据
	GetChildNodeInfo(userid int64, role_code string, status int, account string, page, limit int) ([]UserInfo, error)
	// 查询直接下级用户
	GetChildNodeDirect(userid int64, role_code string, status int) ([]UserInfo, error)
	//check
	CheckChildNodeInfo(userid, childid int64) (bool, error)
	//获取用户密码
	GetUserBase(userid int64) (string, error)
	//  ==权限==
	// 查询可以操作的角色
	//operate: 增:create_role 删delete_role 查select_role  改update_role
	//sale_role  purchase_role  distribute_role  分销相关
	GetOperationRole(userid int64, operate string) ([]UserRole, error)
	// 查询是否有操作权限
	//operate 操作码  roleCode 操作对象的角色码，没有就是""
	GetOperatePermission(userid int64, operate string, roleCode string) (bool, error)

	//==注册==
	// 用户注册
	CreateUser(parentid int64, user *UserBaseInfo, pass string) (int64, error)

	//== 修改用户==
	//修改用户基本信息,这个接口不会修改密码status,parentid,role_code
	UpdateUserBaseInfo(user *UserBaseInfo) error
	//修改用户佣金信息,这个接口不会修改status
	//UpdateUserBalanceInfo(user UserBaseInfo) error
	//修改用户status
	UpdateUserStatus(userid int64, status int) error
	// 修改用户密码
	UpdateUserPass(userid int64, oldpass, newpass string) error
	//管理员帮忙修改密码,operaterid操作的管理员id pass管理员密码
	HelpUpdatePass(account string, operaterid int64,pass, newpass string) error

	//==佣金==
	//增加佣金流水信息。提现金额前端接收的string直接传过来即可
	//CreateMoneyRecord(mr MoneyRecord)(string,error)
	// 申请佣金提现
	ApplyMoneyWithdrew(userid int64, withdrew string) error

	// 佣金清零(上级批量审核佣金提现信息) ，这里没有进行权限校验,请调用接口校验
	// 1. 检查余额 2.生成记录 3.更改用户金额信息
	WithdrewMoneyForChild(userid, operateid int64, withdrew string) error
}

type userDB struct {
	*sql.DB
}

func NewUserDB() (UserDB, error) {
	db, err := datastore.LinkStore.GetDB("master")
	if err != nil {
		return nil, errors.As(err)
	}
	udb := &userDB{db}
	return udb, nil
}

func (db *userDB) Close() {
	db.Close()
}

// 检查用户状态
func (db *userDB) CheckAccount(account string) (bool, error) {
	var status int
	err := db.QueryRow(checkUserByNameSql, account).Scan(
		&status,
	)
	//没数据，系统错误
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.As(UserDataNotFound).As(account)
		} else {
			return false, errors.As(err, account)
		}
	}
	return checkUserStatus(status)
}

//登录校验 返回用户id
func (db *userDB) LoginCheck(account string, pass string) (bool, int64, error) {
	var id int64
	var passFromDB string
	var status int
	err := db.QueryRow(checkLoginSql, account).Scan(
		&id,
		&passFromDB,
		&status,
	)

	//没数据，系统错误
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, errors.As(UserDataNotFound).As(account)
		} else {
			return false, 0, errors.As(err, account)
		}
	}
	// 密码错误
	if !bcrypt.BcryptMatch(pass, passFromDB) {
		return false, 0, errors.As(UserWrongPass)
	}
	var isNormal bool
	isNormal, err = checkUserStatus(status)
	return isNormal, id, errors.As(err, account)
}

//通过id获取用户基本详情信息
func (db *userDB) GetUserBaseInfoById(id int64) (*UserBaseInfo, error) {
	var createTime time.Time
	var updateTime time.Time
	user := UserBaseInfo{}
	err := db.QueryRow(getUserBaseInfoByIdSql, id).Scan(
		&user.UserId,
		&user.RoleCode,
		&user.ParentId,
		//&user.Password,
		&user.Mobile,
		&user.Account,
		&user.RealName,
		&user.IdCard,
		&user.BankCard,
		&user.Address,
		&user.Status,
		&createTime,
		&updateTime,
		&user.ParentName,
		&user.ParentRealName,
		&user.RoleName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.As(UserDataNotFound).As(id)
		} else {
			return nil, errors.As(err, id)
		}
	}

	user.CreateTime = createTime.Format(formatTime)
	user.UpdateTime = updateTime.Format(formatTime)
	user.StatusName = statusName(user.Status)
	return &user, nil
}

// 获取用户佣金数据
func (db *userDB) GetUserBalanceInfoById(userid int64) (*UserBalanceInfo, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.As(err, userid)
	}

	user, err := getUserBalanceInfoById(tx, userid)
	if err != nil {
		tx.Rollback()
		return nil, errors.As(err, userid)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, errors.As(err, userid)
	}

	return user, nil
}

// 获取用户佣金数据事务版
func getUserBalanceInfoById(tx *sql.Tx, userid int64) (*UserBalanceInfo, error) {
	var createTime time.Time
	var updateTime time.Time

	user := UserBalanceInfo{}
	err := tx.QueryRow(getUserBalanceInfoByIdSql, userid).Scan(
		&user.CashNormal,
		&user.ApplyNormal,
		&user.WithdrewNormal,
		&user.TotalNormal,
		//&user.Password,
		&createTime,
		&updateTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// 若不存在，创建一条佣金数据
			_, err := tx.Exec(createUserBalanceInfoSql, userid)
			if err != nil {
				return nil, errors.As(err, userid)
			}
			//重新查
			return getUserBalanceInfoById(tx, userid)
		} else {
			return nil, errors.As(err, userid)
		}
	}

	user.BalanceCreateTime = createTime.Format(formatTime)
	user.BalanceUpdateTime = updateTime.Format(formatTime)
	user.Apply = money.New(float64(user.ApplyNormal) / 10000).Format(2)
	user.Withdrew = money.New(float64(user.WithdrewNormal) / 10000).Format(2)
	user.Cash = money.New(float64(user.CashNormal) / 10000).Format(2)
	user.Total = money.New(float64(user.TotalNormal) / 10000).Format(2)
	return &user, nil
}

//通过账号获取用户基本详情信息
func (db *userDB) GetUserBaseInfoByAccount(account string) (*UserBaseInfo, error) {
	var createTime time.Time
	var updateTime time.Time
	user := UserBaseInfo{}
	err := db.QueryRow(getUserBaseInfoByAccountSql, account).Scan(
		&user.UserId,
		&user.RoleCode,
		&user.ParentId,
		//&user.Password,
		&user.Mobile,
		&user.Account,
		&user.RealName,
		&user.IdCard,
		&user.BankCard,
		&user.Address,
		&user.Status,
		&createTime,
		&updateTime,
		&user.ParentName,
		&user.ParentRealName,
		&user.RoleName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.As(UserDataNotFound).As(account)
		} else {
			return nil, errors.As(err, account)
		}
	}

	user.CreateTime = createTime.Format(formatTime)
	user.UpdateTime = updateTime.Format(formatTime)
	user.StatusName = statusName(user.Status)
	return &user, nil
}

// 获取用户佣金数据
func (db *userDB) GetUserBalanceInfoByAccount(account string) (*UserBalanceInfo, error) {
	var createTime time.Time
	var updateTime time.Time
	user := UserBalanceInfo{}
	err := db.QueryRow(getUserBalanceInfoByAccountSql, account).Scan(
		&user.CashNormal,
		&user.ApplyNormal,
		&user.WithdrewNormal,
		&user.TotalNormal,
		//&user.Password,
		&createTime,
		&updateTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.As(UserDataNotFound).As(account)
		} else {
			return nil, errors.As(err, account)
		}
	}

	user.BalanceCreateTime = createTime.Format(formatTime)
	user.BalanceUpdateTime = updateTime.Format(formatTime)
	user.Apply = money.New(float64(user.ApplyNormal) / 10000).Format(2)
	user.Withdrew = money.New(float64(user.WithdrewNormal) / 10000).Format(2)
	user.Cash = money.New(float64(user.CashNormal) / 10000).Format(2)
	user.Total = money.New(float64(user.TotalNormal) / 10000).Format(2)
	return &user, nil

}

// 获取用户可以操作的角色
func (db *userDB) GetOperationRole(userid int64, operate string) ([]UserRole, error) {
	roleMess, err := db.getAuth(userid, operate)
	if err != nil {
		return nil, errors.As(err, userid, operate)
	}

	return db.getUserRole(roleMess)
}

// 解析用户可以操作的角色
func (db *userDB) getUserRole(roleMess string) ([]UserRole, error) {
	urs := []UserRole{}
	rlist := make([]string, 0)
	json.Unmarshal([]byte(roleMess), &rlist)

	if strings.Index(roleMess, "ALL") >= 0 {
		sqlStr := getAllRoleNameSql
		for _, v := range rlist {
			if v != "ALL" {
				except := fmt.Sprintf(" and role_code!='%s'", v)
				sqlStr = sqlStr + except
			}
		}
		rows, err := db.Query(sqlStr)
		if err != nil {
			if err == sql.ErrNoRows {

				return nil, errors.As(RoleDataNotFound)
			} else {
				return nil, errors.As(err)
			}
		}

		for rows.Next() {
			ur := UserRole{}
			if err := rows.Scan(
				&ur.RoleCode,
				&ur.RoleName,
			); err != nil {
				if err == sql.ErrNoRows {
					return nil, errors.As(RoleDataNotFound)
				} else {
					return nil, errors.As(err)
				}
			}
			urs = append(urs, ur)
		}
	} else {
		for _, v := range rlist {
			ur := UserRole{}
			if err := db.QueryRow(getRoleNameSql, v).
				Scan(&ur.RoleCode, &ur.RoleName); err != nil {
				if err == sql.ErrNoRows {
					return nil, errors.As(RoleDataNotFound)
				} else {
					return nil, errors.As(err)
				}
			}
			urs = append(urs, ur)
		}
	}

	return urs, nil
}

//获得用户的某项权限，返回权限，数组json
func (db *userDB) getAuth(userid int64, authName string) (string, error) {
	var roleMes string
	getRoleByIdSql2 := fmt.Sprintf(getRoleByIdSql, authName)
	if err := db.QueryRow(getRoleByIdSql2, userid).Scan(&roleMes); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.As(RoleDataNotFound).As(userid, authName)
		} else {
			return "", errors.As(err, userid, authName)
		}
	}
	return roleMes, nil
}

//检测用户是否有权限进行某项操作
//operate 操作码  roleCode 操作对象的角色码，没有就是""
func (db *userDB) GetOperatePermission(userid int64, operate string, roleCode string) (bool, error) {
	roleMsg, err := db.getAuth(userid, operate)
	if err != nil {
		return false, errors.As(err, userid, operate, roleCode)
	}

	//查改增删
	switch operate {
	case SelectRole, UpdateRole, DeleteRole, CreateRole, DistributeRole:
		haveRole := strings.Index(roleMsg, roleCode)
		if strings.Index(roleMsg, "ALL") >= 0 {
			if haveRole >= 0 {
				return false, nil
			} else {
				return true, nil
			}
		} else {
			if haveRole < 0 {
				return false, nil
			} else {
				return true, nil
			}
		}
	case SaleRole, PurchaseRole:
		role, err := strconv.ParseBool(roleMsg)
		if err != nil {
			return false, errors.As(err, userid, operate, roleCode)
		}
		return role, nil
	default:
		return false, errors.As(UserWrongAuthParam).As(userid, operate, roleCode)
	}
}

//用户注册
func (db *userDB) CreateUser(parentid int64, info *UserBaseInfo, pass string) (int64, error) {
	var createStatus int
	parent, err := db.GetUserBaseInfoById(parentid)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.As(UserDataNotFound).As(parentid)
		} else {
			return 0, errors.As(err, parentid)
		}
	}

	// 1.创建者是ROOT ADMIN 2.被创建者是消费者  不需要管理员审核。
	if (parent.RoleCode == ADMIN || parent.RoleCode == ROOT) || info.RoleCode == CONSUMER {
		createStatus = Normal
	} else {
		createStatus = Checking
	}

	pass, _ = bcrypt.BcryptHash(pass)
	result, err := db.Exec(createUserBaseInfoSql,
		parentid, pass, info.Mobile, info.Account, info.RealName,
		info.RoleCode, info.IdCard, info.BankCard, createStatus)
	if err != nil {
		return 0, errors.As(err, parentid)
	}

	uid, err := result.LastInsertId()
	if err != nil {
		return 0, errors.As(err, parentid)
	}

	if err := db.createUserBalanceInfo(uid); err != nil {
		return 0, errors.As(err, parentid)
	}

	return uid, nil
}

// 创键一个用户佣金信息表
func (db *userDB) createUserBalanceInfo(userid int64) error {
	_, err := db.Exec(createUserBalanceInfoSql, userid)
	if err != nil {
		return errors.As(err)
	}
	return nil

}

//查询直接上级的基本信息
func (db *userDB) GetParentNodeDirect(userid int64) (*UserBaseInfo, error) {
	var createTime time.Time
	var updateTime time.Time
	user := UserBaseInfo{}
	err := db.QueryRow(getParentNodeInfoSql, userid).Scan(
		&user.UserId,
		&user.RoleCode,
		&user.ParentId,
		//&user.Password,
		&user.Mobile,
		&user.Account,
		&user.RealName,
		&user.IdCard,
		&user.BankCard,
		&user.Address,
		&user.Status,
		&createTime,
		&updateTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.As(UserDataNotFound).As(userid)
		} else {
			return nil, errors.As(err, userid)
		}
	}
	user.CreateTime = createTime.Format(formatTime)
	user.UpdateTime = updateTime.Format(formatTime)
	user.StatusName = statusName(user.Status)
	return &user, nil
}

//查询下级数量,若role_code为空，就是所有角色。若status合法，就筛选status
//若下级用户没有佣金数据，会导致结果为空
func (db *userDB)GetChildNodeNum(userid int64, role_code string, status int, account string, page, limit int)(int,error){
	totalNum := 0
	numStr := getAllChildNodeNumSql
	extraStr := ""
	//数据
	if role_code != "" {
		//若不为空就是查询某个身份的下级
		extraStr = extraStr + fmt.Sprintf(" and tb1.role_code='%s'", role_code)
	}

	if account != "" {
		extraStr = extraStr + " and tb1.account like '%" + account + "%'"
	}

	if isIvalidStatus(status) {
		//若不为空就是查询某个身份的下级
		extraStr = extraStr + fmt.Sprintf(" and tb1.status=%d", status)
	}
	//数据量的sql
	numStr = numStr + extraStr

	if err := db.QueryRow(numStr, userid, userid).Scan(&totalNum); err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				return  0, errors.As(UserDataNotFound).As(userid, role_code, status)
			} else {
				return  0, errors.As(err, userid, role_code, status)
			}
		}
	}
	return totalNum,nil
}

//查询下级,若role_code为空，就是所有角色。若status合法，就筛选status
//若下级用户没有佣金数据，会导致结果为空
func (db *userDB) GetChildNodeInfo(userid int64, role_code string, status int, account string, page, limit int) ([]UserInfo,  error) {
	sqlStr := getAllChildNodeInfoSql

	extraStr := ""
	//数据
	if role_code != "" {
		//若不为空就是查询某个身份的下级
		extraStr = extraStr + fmt.Sprintf(" and tb1.role_code='%s'", role_code)
	}

	if account != "" {
		extraStr = extraStr + " and tb1.account like '%" + account + "%'"
	}

	if isIvalidStatus(status) {
		//若不为空就是查询某个身份的下级
		extraStr = extraStr + fmt.Sprintf(" and tb1.status=%d", status)
	}
	

	extraStr = extraStr + " order by tb1.create_time desc"

	if page > 0 && limit > 0 {
		extraStr = extraStr + fmt.Sprintf(" limit %d offset %d", limit, (page-1)*limit)
	}

	row, err := db.Query(sqlStr+extraStr, userid, userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.As(UserDataNotFound).As(userid, role_code, status)
		} else {
			return nil, errors.As(err, userid, role_code, status)
		}
	}

	uList, err := getUserNodeData(db, row)
	return uList, err
}

//查询直接下级用户
func (db *userDB) GetChildNodeDirect(userid int64, role_code string, status int) ([]UserInfo, error) {
	sqlStr := getChildNodeInfoSql
	//
	if role_code != "" {
		//若不为空就是查询某个身份的下级
		sqlStr = sqlStr + fmt.Sprintf(" and role_code='%s'", role_code)
	}

	if isIvalidStatus(status) {
		//若不为空就是查询某个身份的下级
		sqlStr = sqlStr + fmt.Sprintf(" and status=%d", status)
	}
	sqlStr = sqlStr + " order by create_time desc"

	row, err := db.Query(sqlStr, userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.As(UserDataNotFound).As(userid, role_code, status)
		} else {
			return nil, errors.As(err, userid, role_code, status)
		}
	}
	return getUserNodeData(db, row)
}

//check
func (db *userDB) CheckChildNodeInfo(userid, childid int64) (bool, error) {
	var num int
	err := db.QueryRow(checkChildNodeInfoSql, userid, childid).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, errors.As(err, userid, childid)
		}
	}

	if num == 0 {
		return false, nil
	}
	return true, nil
}

//查询上级,若role_code为空，就是所有角色。若status合法，就筛选status
func (db *userDB) GetParentNodeInfo(userid int64, role_code string, status int) ([]UserInfo, error) {

	sqlStr := getAllParentNodeInfoSql
	//
	if role_code != "" {
		//若不为空就是查询某个身份的下级
		sqlStr = sqlStr + fmt.Sprintf(" and role_code='%s'", role_code)
	}

	if isIvalidStatus(status) {
		//若不为空就是查询某个身份的下级
		sqlStr = sqlStr + fmt.Sprintf(" and status=%d", status)
	}
	sqlStr = sqlStr + " order by create_time desc"

	row, err := db.Query(sqlStr, userid, userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.As(UserDataNotFound).As(userid, role_code, status)
		} else {
			return nil, errors.As(err, userid, role_code, status)
		}
	}
	return getUserNodeData(db, row)
}

func getUserNodeData(db *userDB, row *sql.Rows) ([]UserInfo, error) {
	userList := []UserInfo{}
	for row.Next() {
		user := UserInfo{}

		var createTime time.Time
		var updateTime time.Time
		if err := row.Scan(
			&user.UserId,
			&user.RoleCode,
			&user.ParentId,
			//&user.Password,
			&user.Mobile,
			&user.Account,
			&user.RealName,
			&user.IdCard,
			&user.BankCard,
			&user.Address,
			&user.Status,
			&createTime,
			&updateTime,
			//		&user.RoleName,
			&user.ParentName,
			&user.ParentRealName,
			&user.CashNormal,
			&user.ApplyNormal,
			&user.WithdrewNormal,
			&user.TotalNormal,
			//&user.RoleName,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.As(UserDataNotFound)
			} else {
				return nil, errors.As(err)
			}
		}

		if err := db.QueryRow(getRoleNameSql, user.RoleCode).Scan(&user.RoleCode, &user.RoleName); err != nil {
			if err == sql.ErrNoRows {
				user.RoleName = "无身份用户"
			} else {
				return nil, errors.As(err)
			}
		}
		user.CreateTime = createTime.Format(formatTime)
		user.UpdateTime = updateTime.Format(formatTime)
		user.Apply = money.New(float64(user.ApplyNormal) / 10000).Format(2)
		user.Withdrew = money.New(float64(user.WithdrewNormal) / 10000).Format(2)
		user.Cash = money.New(float64(user.CashNormal) / 10000).Format(2)
		user.Total = money.New(float64(user.TotalNormal) / 10000).Format(2)
		user.StatusName = statusName(user.Status)
		userList = append(userList, user)
	}
	return userList, nil
}

//修改用户基本信息
func (db *userDB) UpdateUserBaseInfo(user *UserBaseInfo) error {
	_, err := db.Exec(updateUserBaseInfoSql,
		user.Mobile, user.Account, user.RealName, user.IdCard, user.BankCard,
		user.Address, user.UserId)
	if err != nil {
		return errors.As(err, user.Account)
	}
	return nil
}

//修改用户status
func (db *userDB) UpdateUserStatus(userid int64, status int) error {
	_, err := db.Exec(updateUserStatusSql,
		status, userid)
	if err != nil {
		return errors.As(err, userid, status)
	}
	return nil
}

// 申请佣金提现
// 1. 检查余额 2.生成记录 3.更改用户金额信息
func (db *userDB) ApplyMoneyWithdrew(userid int64, apply string) error {
	//1. 检查余额
	ub, err := db.GetUserBalanceInfoById(userid)

	if err != nil {
		return errors.As(err, userid)
	}

	applyWant, err := moneyStringToInt(apply)
	if err != nil {
		return errors.As(err, userid)
	}

	if applyWant <= 0 {
		return nil
	}

	if ub.CashNormal < applyWant {
		return errors.As(UserNotEnoughMoney)
	}

	//2.生成记录
	tx, err := db.Begin()
	if err != nil {
		return errors.As(err, userid)
	}

	//目前，状态直接是已到账，不用审核。
	mr := MoneyRecord{
		OrderId: uuid.NewV4().String(),
		UserId:  userid,
		Amount:  strconv.Itoa(applyWant),
		Operate: MoneyRecordApply,
		Status:  MoneyRecordAccount,
		Memo:    "申请佣金提现",
	}
	_, err = createMoneyRecord(tx, mr)
	if err != nil {
		tx.Rollback()
		return errors.As(err, userid)
	}

	//3.更改用户金额信息
	ub.ApplyNormal = ub.ApplyNormal + applyWant
	ub.CashNormal = ub.CashNormal - applyWant
	if err := updateUserBalanceInfo(tx, ub, userid); err != nil {
		tx.Rollback()
		return errors.As(err, userid)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.As(err, userid)
	}
	return nil
}

// 佣金清零(上级批量审核佣金提现信息) ，这里没有进行权限校验,请调用接口校验
// 1. 检查余额 2.生成记录 3.更改用户金额信息
func (db *userDB) WithdrewMoneyForChild(userid, operateid int64, withdrew string) error {
	//1. 检查余额
	ub, err := db.GetUserBalanceInfoById(userid)
	if err != nil {
		return errors.As(err, userid)
	}

	withdrewWant, err := moneyStringToInt(withdrew)
	if err != nil {
		return errors.As(err, userid)
	}

	if ub.ApplyNormal < withdrewWant {
		return errors.As(UserNotEnoughMoney)
	}

	//清零
	if withdrewWant <= 0 {
		withdrewWant = ub.ApplyNormal
	}

	//2.生成记录
	tx, err := db.Begin()
	if err != nil {
		return errors.As(err, userid)
	}

	//目前，状态直接是已到账，不用审核。
	mr := MoneyRecord{
		OrderId:   uuid.NewV4().String(),
		UserId:    userid,
		Amount:    strconv.Itoa(withdrewWant),
		Operate:   MoneyRecordWithdrew,
		OperateId: operateid,
		Status:    MoneyRecordAccount,
		Memo:      "佣金提现申请通过",
	}
	_, err = createMoneyRecord(tx, mr)
	if err != nil {
		tx.Rollback()
		return errors.As(err, userid)
	}

	//3.更改用户金额信息

	ub.ApplyNormal = ub.ApplyNormal - withdrewWant
	ub.WithdrewNormal = ub.WithdrewNormal + withdrewWant

	if err := updateUserBalanceInfo(tx, ub, userid); err != nil {
		tx.Rollback()
		return errors.As(err, userid)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.As(err, userid)
	}
	return nil
}

//增加佣金流水信息。提现金额前端接收的string直接传过来即可
func createMoneyRecord(tx *sql.Tx, mr MoneyRecord) (string, error) {
	if mr.OrderId == "" {
		mr.OrderId = uuid.NewV4().String()
	}

	_, err := tx.Exec(createMoneyRecordSql,
		mr.OrderId,
		mr.UserId,
		mr.OperateId,
		mr.Amount,
		mr.Operate,
		mr.Status,
		mr.Memo,
	)
	if err != nil {
		return "", errors.As(err)
	}

	return mr.OrderId, nil
}

//修改用户佣金信息
func updateUserBalanceInfo(tx *sql.Tx, user *UserBalanceInfo, userid int64) error {
	_, err := tx.Exec(updateUserBalanceInfoSql, user.CashNormal, user.ApplyNormal,
		user.WithdrewNormal, user.TotalNormal, userid)
	if err != nil {
		return errors.As(err)
	}

	return nil
}

// 佣金到账
func accountBalance(tx *sql.Tx, userid int64, cash, memo string) error {
	cashInt, err := moneyStringToInt(cash)
	if err != nil {
		return errors.As(err)
	}
	// 如果是0,直接返回。
	if cashInt <= 0 {
		return nil
	}
	//1. 检查余额
	ub, err := getUserBalanceInfoById(tx, userid)
	if err != nil {
		return errors.As(err)
	}

	//目前，状态直接是已到账MoneyRecordAccount，不用审核。
	mr := MoneyRecord{
		OrderId: uuid.NewV4().String(),
		UserId:  userid,
		Amount:  strconv.Itoa(cashInt),
		Operate: MoneyRecordIN,
		Status:  MoneyRecordAccount,
		Memo:    memo,
	}
	_, err = createMoneyRecord(tx, mr)
	if err != nil {
		tx.Rollback()
		return errors.As(err)
	}

	//2.更改用户金额信息
	ub.TotalNormal = ub.TotalNormal + cashInt
	ub.CashNormal = ub.CashNormal + cashInt

	if err := updateUserBalanceInfo(tx, ub, userid); err != nil {
		tx.Rollback()
		return errors.As(err)
	}

	return nil
}

//修改用户密码
func (db *userDB) UpdateUserPass(userid int64, oldPass, nowPass string) error {
	dbPass, err := db.GetUserBase(userid)
	if err != nil {
		return errors.As(err)
	}

	if !bcrypt.BcryptMatch(oldPass, dbPass) {
		return UserWrongPass.As(userid)
	}

	nowPass, _ = bcrypt.BcryptHash(nowPass)

	_, err = db.Exec(updateUserPassSql, nowPass, userid)
	if err != nil {
		return errors.As(err)
	}

	return nil
}

//管理员帮助修改密码
func (db *userDB) HelpUpdatePass(account string, operaterid int64,pass, newpass string) error {
	//校验管理员密码
	dbPass, err := db.GetUserBase(operaterid)
	if err != nil {
		return errors.As(err)
	}

	if !bcrypt.BcryptMatch(pass, dbPass) {
		return UserWrongPass.As(operaterid)
	}

	operater, err := db.GetUserBaseInfoById(operaterid)
	if err != nil {
		return errors.As(err)
	}

	//校验管理员身份
	if operater.RoleCode == ROOT || operater.RoleCode == ADMIN {

		newpass, _ = bcrypt.BcryptHash(newpass)

		_, err = db.Exec(updateUserPassByAccountSql, newpass, account)
		if err != nil {
			return errors.As(err)
		}
	}else{
		return UserNotEnoughAuth.As(operaterid)
	}

	return nil
}

//获取用户密码
func (db *userDB) GetUserBase(userid int64) (string, error) {
	var base string
	if err := db.QueryRow(getUserBase, userid).Scan(&base); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.As(UserDataNotFound)
		}
		return "", errors.As(err, userid)
	}
	return base, nil
}

//获取父身份
func getParentRole(tx *sql.Tx, roleCode string) (string, error) {
	var parentCode string
	if err := tx.QueryRow(getParentRoleSql, roleCode).Scan(&parentCode); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.As(RoleDataNotFound)
		}
		return "", errors.As(err, roleCode)
	}
	return parentCode, nil
}

type mmsj struct {
	Userid int64
	Pass   string
}

//warning! 群体给pass加密。禁止使用。后果自负。
/*
func TestMMSJ(t *testing.T) {
	dbu.mimashengji()
}
*/
func (db *userDB) mimashengji() {
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	rows, err := tx.Query(mimashengjiQuerySql)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	mmsjs := []mmsj{}
	for rows.Next() {
		m := mmsj{}
		if err := rows.Scan(&m.Userid, &m.Pass); err != nil {
			tx.Rollback()
			panic(err)
		}
		mmsjs = append(mmsjs, m)
	}

	for _, m := range mmsjs {
		pass, _ := bcrypt.BcryptHash(m.Pass)
		_, err := tx.Exec(updateUserPassSql, pass, m.Userid)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		panic(err)
	}

}
