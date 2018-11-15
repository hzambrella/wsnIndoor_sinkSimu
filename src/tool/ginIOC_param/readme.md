# ioc param for gin frame
This package can help you to parse parameter and generate *gin* handler fastly and intelligently.Your code will be shrink.Speed of development will be fase.<br/>
it is not conflict with common way of handler defined.

## function of package
* Get parameter from  query paramter in uri or http post request body automatically.You need not to write code about parse parameter in you controller function ,save lots of time.
* Auto generate gin route handler in target package,named **gen_handler.go**.

 

## how to use

### dependcy
* pkg of frame: "github.com/gin-gonic/gin","tool/gintool"}
* pkg of log "github.com/qiniu/log"
* pkg of error: "tool/errors" 

### support type of parameter and default value when not required:
* buildin type: 

type | default value
---|---
string|""
"int"|  0  
"int8"|0  
"int16"| 0 
"int32"|0 
"int64"|  0 
"uint"| 0 
"uint16"|0 
"uint32"| 0
"uint64"| 0 
"float32|  0
"float64"| 0
"bool"|  false 

* pointer type TODO
* defined in other package TODO
* multipart TODO

###in main.go:
```
import "tool/ginIOC_param"
import "tool/gintool"

func init(){
    //file path of your controller or router 
    //where auto code will be generater.
   ginIOC_param.IOCParam(os.Getenv("PJDIR")+"/routes")
}

func main(){
    r := gintool.Default()
    r.Run(":8080")
}

```
### in controller.go or router.go

* when parameter is required,write paramter name in annotation and formal parameter in function.**ctx must write as "ctx"!**.When param not exist,it will response wrong param message.

```
// format param ctx must write as "ctx"
//@handler GET /t1/:haha?uid_uid&mobile_m
func indexView(ctx *gin.Context, haha int8, uid_uid float64, mobile_m string) {
    //your work
    log.Debug(haha,uid_uid,mobile_m)
    ctx.String(200, "ok")

}
//example of request    http://localhost:8111/t1/12?uid_uid=131&mobile_m=sdad
```
it is equal as common handler style:
```
import "tool/gintool"
func.init(){
      r := gintool.Default()
	r.GET("/t1/:haha", indexViewHandler)  
}

func indexViewHandler(ctx *gin.Context){
    // what you hate and want to omitted
    hahaStr := VarsValue(ctx, "haha")
	pr_haha, err := strconv.ParseInt(hahaStr, 10, 8)
	if err != nil {
		log.Warn(errors.As(err))
		ctx.String(400, "请求参数异常haha")
		return
	}
	haha := int8(pr_haha)
	uid_uidStr := FormValue(ctx, "uid_uid")
	pr_uid_uid, _ := strconv.ParseFloat(uid_uidStr, 64)
	uid_uid := pr_uid_uid
	mobile_mStr := FormValue(ctx, "mobile_m")
	if mobile_mStr == "" { //mobile_m
		log.Warn(errors.As(err))
		ctx.String(400, "请求参数异常mobile_m")
		return
	}
	mobile_m := mobile_mStr
	
	//your actual work
	log.Debug(haha,uid_uid,mobile_m)
    ctx.String(200, "ok")
}

```
So, *it is cool!!!*.These two style is compatible.You can mix use.When too many param to recieve ,it can save your lots of time.


* when parameter isn't required,only write paramter name in formal parameter.These parameter will be pass an default value.You should pay attention to deal with them. 
Like uid_uid in follow

```
//@handler GET /t1/:haha?mobile_m

func indexView(ctx *gin.Context, haha int8, uid_uid float64, mobile_m string) {
	
    log.Debug(haha,uid_uid,mobile_m)
	
    ctx.String(200, "ok")

}
//example of request    http://localhost:8111/t1/12?mobile_m=sdad
//uid_uid will be pass an default value:0.0
```

### in src
go build <br/>
./src //to generate handler file<br/>
ctrl+c<br/>
go build <br/>
./src //start server actually.<br/>

TODO:find way to reduce build
