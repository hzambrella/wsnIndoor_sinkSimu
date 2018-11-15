package routes

//用户登录相关
import (
	"interceptor/auth"
	"model/machine"
	"tool/checkinput"
	"tool/errors"
	"tool/gintool"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
)

const (
	//登录页面
	loginViewPath = "/machine/login/view"
	//检测账号
	CheckAccountPath = "/machine/check/user"
	//检测密码
	CheckUserByPassPath = "/machine/check/pass"
	//注册页面
	registerViewPath = "/machine/register/view"
	//注册
	registerPath = "/machine/register/post"
	//修改密码
	changePassWordViewPath = "/machine/changePassView"
	changePassWordPath     = "/machine/changePass"
	//管理员帮助修改密码
	helpChangePassPath = "/machine/helpChangePass"
	//下线
	logoutPath         = "/machine/logout"
	userDetailViewPath = testTemplatePath
)

func init() {
	r := gintool.Default()
	r.GET(loginViewPath, loginView)
	r.POST(CheckAccountPath, CheckAccount)
	r.POST(CheckUserByPassPath, CheckUserByPass)
	r.GET(registerViewPath, registerView)
	r.POST(registerPath, register)
	r.GET(changePassWordViewPath, changePassWordView)
	r.POST(changePassWordPath, changePassWord)
	r.POST(helpChangePassPath, helpChangePass)
	r.GET(logoutPath, logout)
}

// 登录页面
func loginView(c *gin.Context) {
	//若用户已经登录了，输入这个链接重定向为主页
	uss, err := auth.GetUserSession(c)
	if err != nil {
		if !errors.Equal(err, auth.UserSessNotFound) {
			log.Warn(err)
			c.String(500, sysWrong)
			return
		}
	}

	if uss != nil {
		c.Redirect(302, userDetailViewPath)
		return
	}

	reqPath := auth.GetUserReqPath(c)
	if reqPath == "" {
		reqPath = userDetailViewPath
	}
	c.HTML(200, "machine/login", gin.H{
		"reqPath": reqPath,
	})
}

//检测账号
func CheckAccount(c *gin.Context) {
	db, err := machine.NewUserDB()
	if err != nil {
		log.Warn(err)
		c.String(500, dbWrong)
		return
	}
	//账号
	account := FormValue(c, "account")
	if account == "" {
		log.Warn(paramWrongFormat("account"))
		c.String(400, "请输入用户名")
		return
	}
	var isPass bool
	h := gin.H{}

	//校验用户
	isPass, err = db.CheckAccount(account)

	//无效用户名
	if err != nil {
		h["isValid"] = false
		if errors.Equal(err, machine.UserDataNotFound) {
			h["message"] = machine.NotFoundMes
			h["code"] = UserDataNotFoundCode
		} else if errors.Equal(err, machine.UserFrozen) {
			h["message"] = machine.FrozenMes
			h["code"] = UserFrozenCode
		} else if errors.Equal(err, machine.UserChecking) {
			h["message"] = machine.CheckingMes
			h["code"] = UserCheckingCode
		} else {
			log.Warn(err)
			c.String(500, sysWrong)
			return
		}

		c.JSON(200, h)
		return
	}

	if isPass == true && err == nil {
		h["message"] = "用户有效"
		h["isValid"] = true
		if err := auth.ClearUserReqPath(c); err != nil {
			c.String(500, sysWrong)
			return
		}
		c.JSON(200, h)
	}
}

//校验密码
func CheckUserByPass(c *gin.Context) {
	db, err := machine.NewUserDB()
	if err != nil {
		log.Warn(err)
		c.String(500, dbWrong)
		return
	}
	account := FormValue(c, "account")
	if account == "" {
		log.Warn(paramWrongFormat("account"))
		c.String(400, "请输入用户名")
		return
	}

	passWord := FormValue(c, "passWord")
	if passWord == "" {
		log.Warn(paramWrongFormat("passWord"))
		c.String(400, "请输入密码")
		return
	}

	//校验用户
	isPass, uid, err := db.LoginCheck(account, passWord)
	h := gin.H{}

	if err != nil {
		//无效用户名
		h["isValid"] = false
		//无效密码
		h["wrongPass"] = false

		if errors.Equal(err, machine.UserDataNotFound) {
			h["message"] = machine.NotFoundMes
		} else if errors.Equal(err, machine.UserFrozen) {
			h["message"] = machine.FrozenMes
		} else if errors.Equal(err, machine.UserChecking) {
			h["message"] = machine.CheckingMes
		} else if errors.Equal(err, machine.UserWrongPass) {
			h["message"] = machine.WrongPassMes
			h["isValid"] = true
			h["wrongPass"] = true
		} else {
			log.Warn(err)
			c.String(500, sysWrong)
			return
		}

		c.JSON(200, h)
		return
	}

	if isPass == true && err == nil {
		if err := auth.SetUserSession(uid, c); err != nil {
			log.Warn(err)
			c.String(500, sysWrong)
			return
		}

		h["message"] = "密码正确"
		h["isValid"] = true
		c.JSON(200, h)
		return
	}
}

// 注册页面
func registerView(c *gin.Context) {
	if !auth.Auth(c) {
		return
	}

	uss, err := auth.GetUserSession(c)
	if err != nil {
		log.Warn(err)
		c.String(500, sysWrong)
		return
	}

	h := gin.H{}
	db, err := machine.NewUserDB()
	if err != nil {
		log.Warn(err)
		c.String(500, dbWrong)
		return
	}

	ussRoleList, err := db.GetOperationRole(uss.Uid, machine.CreateRole)
	if err != nil {
		log.Warn(err)
		if errors.Equal(err, machine.RoleDataNotFound) {
			c.String(500, "数据库角色权限信息异常")
			return
		} else {
			c.String(500, dbWrong)
			return
		}
	}

	resultList := make([]map[string]interface{}, 0)
	for _, v := range ussRoleList {
		hh := make(map[string]interface{}, 0)
		hh["roleCode"] = v.RoleCode
		hh["roleName"] = v.RoleName
		resultList = append(resultList, hh)
	}
	h["resultList"] = resultList

	c.HTML(200, "machine/register", h)
}

//注册
func register(c *gin.Context) {
	if !auth.Auth(c) {
		return
	}

	db, err := machine.NewUserDB()
	if err != nil {
		log.Warn(err)
		c.String(500, dbWrong)
		return
	}
	// parentid
	uss, err := auth.GetUserSession(c)
	if err != nil {
		log.Warn(err)
		c.String(500, sysWrong)
		return
	}
	/*
		parentidStr := FormValue(c, "parentid")
		if parentidStr == "" {
			log.Warn(paramWrongFormat("parentid"))
			c.String(400, paramWrongFormat("parentid"))
			return
		}

		parentid, err := strconv.ParseInt(parentidStr, 10, 64)
		if err != nil {
			log.Warn(err)
			c.String(500, sysWrong)
			return
		}
	*/

	parentid := uss.Uid

	//密码
	passWord := FormValue(c, "passWord")
	if passWord == "" {
		log.Warn(paramWrongFormat("passWord"))
		c.String(400, paramWrongFormat("密码"))
		return
	}
	//用户名
	account := FormValue(c, "account")
	if account == "" {
		log.Warn(paramWrongFormat("account"))
		c.String(400, paramWrongFormat("用户名"))
		return
	}

	//真实名
	realName := FormValue(c, "realName")
	/*
		if realName == "" {
			log.Warn(paramWrongFormat("realName"))
			c.String(400, paramWrongFormat("真实姓名"))
			return
		}
	*/
	//角色
	roleCode := FormValue(c, "roleCode")
	if roleCode == "" {
		log.Warn(paramWrongFormat("roleCode"))
		c.String(400, paramWrongFormat("角色名"))
		return
	}
	//身份证号
	iDCard := FormValue(c, "idCard")
	/*
		if iDCard == "" {
			log.Warn(paramWrongFormat("idCard"))
			c.String(400, paramWrongFormat("身份证号"))
			return
		}
	*/
	// 银行卡
	bankCard := FormValue(c, "bankCard")
	/*
		if bankCard == "" {
			log.Warn(paramWrongFormat("bankCard"))
			c.String(400, paramWrongFormat("银行卡号"))
			return
		}
	*/

	userR := &machine.UserBaseInfo{
		Account:  account,
		RealName: realName,
		IdCard:   iDCard,
		BankCard: bankCard,
	}

	userR.RoleCode = roleCode
	//检查用户名是否存在
	uisExist, err := db.GetUserBaseInfoByAccount(account)
	if err != nil {
		if !errors.Equal(err, machine.UserDataNotFound) {
			log.Warn(err)
			c.String(500, dbWrong)
			return
		}
	}
	if uisExist != nil {
		c.String(400, "用户名已经被使用")
		return
	}

	// 手机号
	/*
		mobile := FormValue(c, "mobile")
		if bankCard == "" {
			log.Warn(paramWrongFormat("mobile"))
			c.String(400, paramWrongFormat("手机号"))
			return
		}
	*/

	//校验账号
	//若为消费者，账号必须是11位手机号
	if roleCode == machine.CONSUMER {
		if !checkinput.IsPhone(account) {
			c.String(400, "消费者的账号必须为合法的11位手机号")
			return
		}

		// 手机号
		userR.Mobile = account
	}

	uid, err := db.CreateUser(parentid, userR, passWord)
	if err != nil {
		log.Warn(err)
		c.String(500, sysWrong)
		return
	}

	c.JSON(200, gin.H{
		"userid": uid,
	})
}

//修改密码
func changePassWordView(c *gin.Context) {
	if !auth.Auth(c) {
		return
	}

	c.HTML(200, "machine/changePass", gin.H{})
}

func changePassWord(c *gin.Context) {
	if !auth.Auth(c) {
		return
	}

	db, err := machine.NewUserDB()
	if err != nil {
		log.Warn(err)
		c.String(500, dbWrong)
		return
	}
	// parentid
	uss, err := auth.GetUserSession(c)
	if err != nil {
		log.Warn(err)
		c.String(500, sysWrong)
		return
	}

	oldPass := FormValue(c, "old")
	if oldPass == "" {
		c.String(400, paramWrongFormat("原始密码为空"))
		return
	}

	newPass := FormValue(c, "new")
	if newPass == "" {
		c.String(400, paramWrongFormat("新密码为空"))
		return
	}

	if err := db.UpdateUserPass(uss.Uid, oldPass, newPass); err != nil {
		if errors.Equal(err, machine.UserWrongPass) {
			c.String(400, machine.WrongPassMes)
			return
		} else {
			log.Warn(err)
			c.String(500, dbWrong)
			return
		}
	}

	c.String(200, success)

}

func helpChangePass(c *gin.Context) {
	if !auth.Auth(c) {
		return
	}

	db, err := machine.NewUserDB()
	if err != nil {
		log.Warn(err)
		c.String(500, dbWrong)
		return
	}
	// parentid
	uss, err := auth.GetUserSession(c)
	if err != nil {
		log.Warn(err)
		c.String(500, sysWrong)
		return
	}

	pass := FormValue(c, "pass")
	if pass == "" {
		c.String(400, paramWrongFormat("管理员密码为空"))
		return
	}

	newPass := FormValue(c, "new")
	if newPass == "" {
		c.String(400, paramWrongFormat("新密码为空"))
		return
	}

	account := FormValue(c, "account")
	if newPass == "" {
		c.String(400, paramWrongFormat("被帮助者为空"))
		return
	}

	if err := db.HelpUpdatePass(account, uss.Uid, pass, newPass); err != nil {
		if errors.Equal(err, machine.UserWrongPass) {
			c.String(400, machine.WrongPassMes)
			return
		} else if errors.Equal(err, machine.UserNotEnoughAuth) {
			c.String(400, "您的权限不足")
			return
		} else {
			log.Warn(err)
			c.String(500, dbWrong)
			return
		}

	}

	c.String(200, success)

}

func logout(c *gin.Context) {
	if err := auth.DeleteUserSession(c); err != nil {
		log.Warn(err)
		c.String(500, sysWrong)
		return
	}
	//auth.GetUserSession(c)

	//清除cookie
	/*
		machineCookie,err:=c.Request.Cookie("machine_user_cookie")
		if err!=nil{
			log.Debug(err)
		}

		if machineCookie!=nil{
			machineCookie.MaxAge=-1
			c.SetCookie(
			machineCookie.Name,
			machineCookie.Value,
			machineCookie.MaxAge,
			machineCookie.Path,
			machineCookie.Domain,
			machineCookie.Secure,
			machineCookie.HttpOnly,
			)
		}
	*/
	//c.String(200,"ok")
	c.Redirect(302, "/machine/login/view")
}
