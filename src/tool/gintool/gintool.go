package gintool

import (

	"html/template"
	"sync"

	"tool/inicfg"
	"tool/errors"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

//var temp *template.Template
var (
	// Global instance
	defaultE     *gin.Engine
	defaultELock = sync.Mutex{}
)

//设置路由，设置过滤器
func Default() *gin.Engine {
	defaultELock.Lock()
	defer defaultELock.Unlock()

	if defaultE == nil {
		app, err := inicfg.Getcfg().GetSection("app")
		if err != nil {
			panic(err)
		}

		appName:=app["name"]
		if appName==""{
			panic(errors.New("panic:appName is nil!"))
		}

		engine := gin.Default()
		/*
		   if temp!=nil{
		       engine.SetHTMLTemplate(temp)
		   }
		*/
		//engine.Use(gin.Logger(), gin.Recovery())
		defaultE = engine

		//为啥这里的过滤器才生效？？？谁能告诉我！！！
/*
		// redis session

		mapEtc, err := inicfg.Getcfg().GetSection("redis")
		if err != nil {
			panic(err)
		}
		store_redis,err:= sessions.NewRedisStore(4, "tcp", mapEtc["web_redis_uri"], "",[]byte(appName+"_user_redis"))
		if err!=nil{
			panic(err)
		}
		defaultE.Use(sessions.Sessions(appName+"_user_redis", store_redis))
		*/

		store_cookie:=sessions.NewCookieStore([]byte(appName+"_user_cookie"))
		//不能乱创 cookie。防止莫名的bug
		store_cookie.Options(sessions.Options{MaxAge: -1, Path: "/"+appName,HttpOnly:true})
		defaultE.Use(sessions.Sessions(appName+"_user_cookie", store_cookie))
		//defaultE.Use(auth.Auth)
		// cookie session

		//404 seo处理
		//可以采用r.NoRoute()
		// defaultE.Use(func(c *gin.Context){
		// 	c.Next()
		// 	if c.Writer.Status()==404{
		// 		c.HTML(200,"pageNotfound",gin.H{
		// 			"title":"纳尼，这个页面已经躲猫猫了",
		// 			"message":"点这个的链接试试",
		// 			"appName":appName,
		// 		})
		// 	}
		// 	c.Abort()
		// })
	}
	return defaultE
}

func SetHTMLTemplate(templ *template.Template) {
	defaultE.HTMLRender = render.HTMLProduction{Template: templ}
}
