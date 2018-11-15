package auth

import (
	"time"
	"github.com/satori/go.uuid"
	"encoding/json"
	"tool/errors"
	
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
	//"tool/redis"
)

//鉴权拦截器
var (
	UserSessNotFound = errors.New("user not found")
	// session保存时间
	useSessionMagAge=1*60 * 60
)


//用户会话信息
type UserSession struct {
	//uid_auth
	//AuthType string `json:"auth_type"`
	Uid int64 `json:"uid"`
	Random string `json:"random"`
	LastTime string `json:"time"` //会话时间
}

const (
	// 用户会话在session的key
	UserSessionKey = "machine_user"
	// 用户访问请求链接在session的key
	ReqPathSessionKey="machine_req_path"
	//登录页面链接,在routes/login.go
	loginViewPath  = "/machine/login/view"
)

//鉴权
//拦截器千万别用r.Use()，放到handler的第一行。
func Auth(c *gin.Context)bool {
	/*
	needAuth, err := NeedInterceptor(c.Request.URL.Path)
	if err != nil {
		c.AbortWithError(500, err)
		c.String(500, err.Error())
		return
	}
	if !needAuth {
		c.Next()
		return
	}
*/
	
	userSess, err := GetUserSession(c)
	if userSess!=nil{
		log.Debug("auth:",userSess.Uid ," last time:",userSess.LastTime)
	}else{
		log.Debug("auth: not auth")
	}

	if err != nil {
		if err != UserSessNotFound {
			c.AbortWithError(500, err)
			c.String(500, err.Error())
			return false
		}else{
			url:=c.Request.URL.RequestURI()
			log.Debug(url)
			if err:=SaveUserReqPath(url,c);err!=nil{
				c.AbortWithError(500, err)
				c.String(500, err.Error())
				return false
			}
			c.Redirect(302, "/machine/login/view")
			return false
		}
	}

	//重置cookie
	if err := setUserSession(userSess, c); err != nil {
		c.AbortWithError(500, err)
		c.String(500, err.Error())
		return false
	}

	return true
}

/*
func isLogin(c *gin.Context)bool{
	)
	JSONsess.Get("user_login")
}
*/
// 设置会话
func SetUserSession(uid int64, c *gin.Context) error {
	//删除老会话
	if err:=DeleteUserSession(c);err!=nil{
		return err
	}

	//UserSession
	uss := &UserSession{Uid: uid,Random:uuid.NewV4().String(),LastTime:time.Now().Format("2006-01-02 15:04:05")}
	return setUserSession(uss, c)
}

func setUserSession(user *UserSession, c *gin.Context) error {
	sess := sessions.Default(c)
	userAuth, err := json.Marshal(user)
	if err != nil {
		return errors.As(err)
	}
	sess.Options(sessions.Options{MaxAge: useSessionMagAge, Path: "/machine",HttpOnly:true})
	sess.Set(UserSessionKey, userAuth)
	if err:=sess.Save();err!=nil{
		return errors.As(err)
	}
	return nil
}

// 获取会话
func GetUserSession(c *gin.Context) (*UserSession, error) {
	sess := sessions.Default(c)
	uSessByte := sess.Get(UserSessionKey)
	if uSessByte == nil {
		return nil, UserSessNotFound
	}
	userSess := &UserSession{}
	if err := json.Unmarshal(uSessByte.([]byte), userSess); err != nil {
		return nil, errors.As(err)
	}
	return userSess, nil
}

//删除会话，退出登录
func DeleteUserSession(c *gin.Context)error{
	sess := sessions.Default(c)
	sess.Options(sessions.Options{MaxAge: -1, Path: "/machine",HttpOnly:true})
	sess.Delete(UserSessionKey)
	if err:=sess.Save();err!=nil{
		return errors.As(err)
	}
	return nil
}
/*
// 请求是否需要拦截鉴权，在$ETCDIR/intercepotor.ini设置不需拦截的。
func NeedInterceptor(paths string) (bool, error) {
	_, err := os.Stat(os.Getenv("ETCDIR") + "/interceptor.ini")
	if err != nil {
		if os.IsNotExist(err) {
			return true, errors.New(os.Getenv("ETCDIR") + "/interceptor.ini not found")
		}
	}

	mapExclude, err := inicfg.Getcfg().GetSection("intercetor_exclude")
	if err != nil {
		return true, errors.As(err)
	}

	exclude, ok := mapExclude["exclude_path"]
	if !ok {
		return true, errors.New("wrong config about interceptor,please check in " + os.Getenv("ETCDIR") + "/intercepotor.ini")
	}

	result := strings.Index(exclude, paths)
	if result < 0 {
		return true, nil
	}

	return false, nil
}
*/
//用户第一次请求的链接
//清除
func ClearUserReqPath(c *gin.Context)error{
	sess := sessions.Default(c)
	sess.Delete(ReqPathSessionKey)
	if err:=sess.Save();err!=nil{
		return errors.As(err)
	}
	return nil
}

//获取
func GetUserReqPath(c *gin.Context)string{
	sess := sessions.Default(c)
	reqPath:=sess.Get(ReqPathSessionKey)
	if reqPath==nil{
		return ""
	}else{
		return reqPath.(string)
	}

}

//保存
func SaveUserReqPath(path string, c *gin.Context) error {
	sess := sessions.Default(c)
	sess.Options(sessions.Options{MaxAge: useSessionMagAge, Path: "/machine"})
	sess.Set(ReqPathSessionKey, path)
	if err:=sess.Save();err!=nil{
		return errors.As(err)
	}
	return nil
}

func (u *UserSession) GetUid() int64 {
	return u.Uid
}
