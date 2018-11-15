package routes

import (
	//"math/rand"
	//"time"

	//"github.com/gin-gonic/contrib/sessions"
	//"interceptor/auth"
	"tool/gintool"
	//"tool/captcha"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
)

//框架测试
const (
	testJSONPath     = "/machine/json"
	testTemplatePath = "/machine/template"
	IndexViewPath    = "/machine/testdb"
	VueViewPath      = "/machine/vue"

	//验证码测试
	getCaptchaPath   = "/machine/getCaptcha"
	checkCaptchaPath = "/machine/checkCaptcha"
)

func init() {
	r := gintool.Default()
	r.GET(testTemplatePath, testTemplate)
	r.GET(testJSONPath, testJSON)
	r.GET(VueViewPath, VueView)
	//r.GET(IndexViewPath, indexView)
	//r.GET(getCaptchaPath,getCaptcha)
	//r.POST(checkCaptchaPath,checkCaptcha)
}

// 主页面
func indexView(ctx *gin.Context, haha int8, uid_uid float64, mobile_m string) {
	log.Debug(haha, uid_uid, mobile_m)
	ctx.String(200, "ok")
}

func testJSON(c *gin.Context) {
	c.JSON(200, "test ok")

}

func testTemplate(c *gin.Context) {
	/*
		if !auth.Auth(c) {
			return
		}
	*/
	/*
		userSess, err := auth.GetUserSession(c)
		if err != nil {
			c.String(500, err.Error())
		}
	*/

	//	log.Debug(base64img)
	c.HTML(200, "machine/test", gin.H{
		"map": "template is 欧克",
	})
}

func VueView(c *gin.Context) {
	c.HTML(200, "machine/vue", nil)
}

/*
//做实验：验证码功能
func getCaptcha(c *gin.Context) {
	var idKey string
	var base64img string
	//产生随机数
	rand.Seed(time.Now().UnixNano())
	randNum:=rand.Intn(2)
	log.Debug(randNum)
	if randNum==1{
		idKey,base64img=captcha.CreateDigitsCode()
	}else{
		idKey,base64img=captcha.CreateDataCode()
	}

	h:=gin.H{}
	h["idkey"]=idKey
	h["base64img"]=base64img
	c.JSON(200,h)
}

func checkCaptcha(c *gin.Context) {
	idkey:=FormValue(c,"idkey")
	value:=FormValue(c,"value")

	isRight:=captcha.VerfiyCaptcha(idkey,value)
	if isRight{
		c.String(200,"success")
	}else{
		c.String(400,"wrong captcha")
	}
}
*/
