package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
	"strconv"
	//"os"
	"github.com/qiniu/log"
	_ "routes"
	"tool/errors"
	"tool/ginIOC_param"
	"tool/gintool"
	"tool/inicfg"
)

/*
func init() {
	cfg, err := inicfg.Newcfg(os.Getenv("ETCDIR"))
	if err != nil {
		panic(err)

	}
	//test
	_, err = cfg.GetSection("master")
	if err != nil {
		panic(err)
	}
}
*/
func init() {
	ginIOC_param.IOCParam(os.Getenv("PJDIR") + "/routes")
	mapWeb, err := inicfg.Getcfg().GetSection("engine/log")
	if err != nil {
		panic(err)
	}
	logLevel, err := strconv.Atoi(mapWeb["level"])
	if err != nil {
		panic(err)
	}
	log.SetOutputLevel(logLevel)
}

func main() {

	r := gintool.Default()
	gin.SetMode(gin.ReleaseMode)
	//gin.SetMode(gin.ReleaseMode)
	tpl := template.Must(template.New("machine").ParseGlob("./public/html/*.html"))
	gintool.SetHTMLTemplate(tpl)

	r.Static("/public/js", "./public/js")
	r.Static("/public/css", "./public/css")
	r.Static("/public/images", "./public/images")

	app, err := inicfg.Getcfg().GetSection("app")
	if err != nil {
		panic(err)
	}
	appName := app["name"]
	if appName == "" {
		panic(errors.New("panic:appName is nil!"))
	}
	//404
	r.NoRoute(func(c *gin.Context) {
		if c.Writer.Status() == 404 {
			c.HTML(200, "pageNotfound", gin.H{
				"title":   "纳尼，这个页面已经躲猫猫了",
				"message": "点这个的链接试试",
				"appName": appName,
			})
		}
	})

	mapWeb, err := inicfg.Getcfg().GetSection("web/server/uri")
	if err != nil {
		panic(err)
	}

	log.Info("start at:", mapWeb["server_port"])
	r.Run(mapWeb["server_port"])
}
