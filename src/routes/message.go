package routes

//信息
const (
	//成功信息
	success = "操作成功"
	//数据库错误
	dbWrong = "数据库异常"
	//参数错误
	paramWrong = "请求参数异常"
	//系统错误
	sysWrong = "系统服务异常"
	//配置文件错误
	etcWrong = "系统配置异常"
)

//参数错误信息+错误的参数
func paramWrongFormat(varname string) string {
	return paramWrong + ":" + varname
}

//用户信息状态码
const (
	//账号正常
	UserNormalCode = 1000
	//账号被冻结
	UserFrozenCode = 1001
	//账号待审核
	UserCheckingCode = 1002
	//账号不存在
	UserDataNotFoundCode = 1004
	//用户系统异常
	UserSysWrongCode = 1005
)

//用户角色操作信息状态码
const (
	//操作成功
	SuccessCode = 2000
	//没有权限
	NoRightCode = 2001
)
