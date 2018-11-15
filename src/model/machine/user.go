package machine

import (
	"strconv"
	"tool/errors"
	"tool/money"
)

//用户异常状态的提示语
const (
	FrozenMes    = "您的账户已被冻结"
	CheckingMes  = "您的账户正在审核"
	NotFoundMes  = "用户不存在"
	WrongPassMes = "密码错误"
	WrongAuthMes = "权限校验参数不存在"
)

type UserInfo struct {
	UserBaseInfo
	UserBalanceInfo
}

//用户基本信息
type UserBaseInfo struct {
	UserId int64 `json:"user_id"`
	UserRole
	ParentId       int    `json:"parent_id"`
	ParentName     string `json:"parent_name"`      //用户名
	ParentRealName string `json:"parent_real_name"` //真实名
	//Password   string `json:"password"`
	Mobile     string `json:"mobile"`
	Account    string `json:"account"`
	RealName   string `json:"real_name"`
	IdCard     string `json:"id_card"`
	BankCard   string `json:"bank_card"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	StatusName string `json:"status_name"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

//用户金额信息
type UserBalanceInfo struct {
	//UserId int64
	// 给前端的显示的数据
	//可提现佣金
	Cash string
	//拥金余额(压款+可提现金额)
	Apply string
	//已提现佣金
	Withdrew string
	//佣金总额
	Total string
	//数据库原始数据,厘为单位，用于计算。
	CashNormal        int
	ApplyNormal       int
	WithdrewNormal    int
	TotalNormal       int
	BalanceCreateTime string
	//更新时间
	BalanceUpdateTime string
}

//佣金流水记录
type MoneyRecord struct {
	OrderId    string //订单id
	UserId     int64  //用户ID
	OperateId  int64  // 经办人id,提现时审核的人记录于此
	Amount     string //金额
	Operate    string //操作 APPLY代表申请提现，WITHDREW代表已提现，IN代表佣金进账。
	CreateTime string //创建时间
	UpdateTime string //更新时间
	Status     int    //0,审核中; 1，已到账，2，已取消
	Memo       string //备注。获取时，将机器订单号填写于此
}

//佣金流水状态
const (
	MoneyRecordChecking = iota //0,审核中
	MoneyRecordAccount         //1，已到账
	MoneyRecordCancel          //2，已取消
)

//佣金记录操作
//APPLY代表申请提现，WITHDREW代表已提现，IN代表佣金进账。
const (
	MoneyRecordApply    = "APPLY"
	MoneyRecordWithdrew = "WITHDREW "
	MoneyRecordIN       = "IN"
)

//用户角色
type UserRole struct {
	RoleCode string `json:"role_code"`
	RoleName string `json:"role_name"`
}

//用户状态
const (
	Normal = iota
	//冻结
	Frozen
	//审核中
	Checking
)

//用户状态中文名
const (
	NormalName = "已审核"
	//冻结
	FrozennName = "冻结"
	//审核中
	CheckingName = "待审核"
)

//用户角色
const (
	ROOT     = "ROOT"     //系统管理员
	ADMIN    = "ADMIN"    //平台管理员
	CITY     = "CITY"     //城市运营中心
	MANAGER  = "MANAGER"  //管理中心
	SERVER   = "SERVER"   //服务运营中心
	CKZX     = "CKZX"     //创客中心
	CB       = "CB"       //消费商
	CONSUMER = "CONSUMER" //消费者
)

//用户操作
const (
	CreateRole     = "create_role"
	UpdateRole     = "update_role"
	DeleteRole     = "delete_role"
	SelectRole     = "select_role"
	DistributeRole = "distribute_role"
	PurchaseRole   = "purchase_role" //买实物操作。有这个权限的人才能购买实物。
	SaleRole       = "sale_role"     //卖实物操作。有这个权限的人才能将下单地址发送出去。
)

//用户操作名字
const (
	CreateRoleName     = "新建"
	UpdateRoleName     = "更改"
	DeleteRoleName     = "删除"
	SelectRoleName     = "选择"
	DistributeRoleName = "分配机器码"
	PurchaseRoleName   = "购买实物" //买实物操作。有这个权限的人才能购买实物。
	SaleRoleName       = "销售实物" //卖实物操作。有这个权限的人才能将下单地址发送出去。
)

//校验用户状态
func checkUserStatus(status int) (bool, error) {
	switch status {
	case Frozen:
		return false, UserFrozen
	case Checking:
		return false, UserChecking
	}
	return true, nil
}

//是否是有效的用户状态
func isIvalidStatus(status int) bool {
	if status == Normal || status == Frozen || status == Checking {
		return true
	}
	return false
}

//用户状态对应的状态名字
func statusName(status int) string {
	switch status {
	case Normal:
		return NormalName
	case Frozen:
		return FrozennName
	case Checking:
		return CheckingName
	}
	return "未知状态"
}

//用户权限不足提示
//operation 操作, from 角色名, to 角色名
func NoPermissionInfo(operation, from, to string) string {
	var name string
	switch operation {
	case CreateRole:
		name = CreateRoleName
	case UpdateRole:
		name = UpdateRoleName
	case DeleteRole:
		name = DeleteRoleName
	case SelectRole:
		name = SelectRoleName
	case DistributeRole:
		name = DistributeRoleName
	case PurchaseRole:
		name = PurchaseRoleName
	case SaleRole:
		name = SaleRoleName
	default:
		return "权限信息错误"
	}
	return "您没有权限:" + from + "->" + name + "->" + to
}

func moneyStringToInt(money string) (int, error) {
	amount, err := strconv.ParseFloat(money, 64)
	if err != nil {
		return -1, errors.As(err)
	}

	amount2, err := strconv.Atoi(strconv.FormatFloat(amount*10000, 'f', -1, 64))
	if err != nil {
		return -1, errors.As(err)
	}
	return amount2, nil
}

func moneyAdd(balance ...string) (string, error) {
	sumMoney := money.New(0)
	for _, v := range balance {
		amount, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return "-1", errors.As(err)
		}
		sumMoney.Add(money.New(amount))
	}
	sum := sumMoney.String()
	return sum, nil
}
