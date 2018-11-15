package ginIOC_param

import (
	"fmt"
	"go/ast"
	//"log"
	"net/url"
	"regexp"
	"strings"
)

var builtinTypes = map[string]struct{}{
	"string":  struct{}{},
	"int":     struct{}{},
	"int8":    struct{}{},
	"int16":   struct{}{},
	"int32":   struct{}{},
	"int64":   struct{}{},
	"uint":    struct{}{},
	"uint16":  struct{}{},
	"uint32":  struct{}{},
	"uint64":  struct{}{},
	"float32": struct{}{},
	"float64": struct{}{},
	"bool":    struct{}{},
}

var paramWrong string = "\"请求参数异常%s\""

var dealErr string = "		log.Warn(errors.As(err));ctx.String(400," + paramWrong + ");return;"

var (
	varsRegexp       = regexp.MustCompile(`(:[^/]+)`)
	handlerTagRegexp = regexp.MustCompile(`@handler (.+) (.+)`)
)

func (g *generator) pickupImports(x string, imports map[string]string) {
	for pkg, alias := range imports {
		if alias == x {
			g.importPkg(pkg, alias)
		} else if strings.HasSuffix(pkg, x) {
			g.importPkg(pkg, alias)
		}
	}
}

func (g *generator) generateHandler(gt *genTag, imports map[string]string) {
	var (
		names    []string
		priority = make(map[string]string)
	)

	// 如果@handler标记的不是函数，直接返回
	fn, ok := gt.Decl.(*ast.FuncDecl)
	if !ok || fn.Recv != nil {
		panic("invalid @handler tag" + fn.Name.Name)
	}

	// 格式：!METHOD path?querystring
	// !表示不需要自动生成router代码，可选
	// METHOD在HTTP方法GET、POST、PUT、DELETE、OPTION中选择一个
	// path代表请求的地址，path可以携带‘命名’参数，以:为标识符，如/users/:user
	// querystring代表这个请求可能携带的表单参数（用于确定参数获取的方式）
	// 例如：POST /test/:id/info?uid&mobile
	// !未设置，METHOD=POST，PATH=/test/:id/info，querystring=uid&mobile
	matchs := handlerTagRegexp.FindStringSubmatch(gt.Tag)
	if len(matchs) != 3 {
		panic("invalid @handler tag.")
	}

	u, err := url.ParseRequestURI(matchs[2])
	if err != nil {
		panic("invalid request uri")
	}

	// 未带有!的情况，自动生成router注册代码
	// 例如：genRouter.Post("/test", testHandler)
	// 其中test是表示有@handler函数的名称，/test是上面获得的path
	if !strings.HasPrefix(matchs[1], "!") {
		if bytes := g.initBuffer.Bytes(); len(bytes) == 0 {
			g.pInit("r:= gintool.Default()")
		}
		g.pInit("r.", matchs[1], `("`, u.Path, `", `, fmt.Sprint(fn.Name), "Handler)")
	}

	// 找出path中的‘命名’参数
	// ‘命名’参数的优先级比FORM表单优先级低，如果后面出现相同参数将会被覆盖
	vars := varsRegexp.FindAllStringSubmatch(u.Path, -1)

	for _, v := range vars {
		if len(v) > 0 {
			priority[v[0][1:]] = "ginIOC_param.Vars"
		}
	}
	// 找出FORM表单中的参数
	for k, _ := range u.Query() {
		priority[k] = "ginIOC_param.Form"
	}

	g.p("func ", fn.Name.Name, "Handler(ctx *gin.Context) {")
	for _, parameter := range fn.Type.Params.List {
		if parameter.Names[0].Name == "ctx" {
			continue
		}
		// 除了builtin类型，其他类型都从请求的body中进行反序列化
		for _, name := range parameter.Names {
			g.importPkg("strconv","")
			g.importPkg(thisPkgName,"")
			switch t := parameter.Type.(type) {
			case *ast.Ident:
				// builtin类型和当前包定义的类型
				if _, ok := builtinTypes[t.Name]; ok {
					/*
						var method string
						// 首字母大写，用于拼接求值函数
						ttype := strings.Title(t.Name)
						if p, ok := priority[name.Name]; ok {
							method = p + ttype
						} else {
							// 如果没有找到相应的参数，那么从协议里面获取
							// 协议必须实现builtin类型参数的获取方法
							method = "Get" + ttype
						}
						g.p(name.Name, ", err := ctx.", method, `("`, toSnake(name.Name, gt.Snake), `")`)
						g.p("if err != nil { log.Warn(errors.As(err)); return }")
					*/
					methodFront, required := priority[name.Name]
					if methodFront == "" || !required {
						methodFront = "Form"
					}

					g.p(name.Name+"Str:=", methodFront, `Value(ctx,"`, toSnake(name.Name, gt.Snake), `")`)
					var nextParse []string
					if required {
						g.importPkg(logPkgName,"")
						g.importPkg(errorsPkgName,"")
						nextParse = builtinParseMethod[t.Name][0]
						nextParse[0] = fmt.Sprintf(nextParse[0], name.Name, name.Name)
						nextParse[1] = fmt.Sprintf(nextParse[1], name.Name)
						nextParse[3] = fmt.Sprintf(nextParse[3], name.Name, name.Name)

					} else {
						nextParse = builtinParseMethod[t.Name][1]
						nextParse[0] = fmt.Sprintf(nextParse[0], name.Name, name.Name)
						nextParse[1] = fmt.Sprintf(nextParse[1], name.Name, name.Name)
					}
					for _, v := range nextParse {
						g.p(v)
					}

				} else {
					panic(fmt.Sprintf("not supported type: %#v\n :(param %s)", t, name.Name))
					/*
						// 当前包自定义的类型
						g.p("var ", name.Name, " ", t.Name)
						g.p("if err := ctx.From", gt.Format, "(&", toSnake(name.Name, gt.Snake), "); err != nil { log.Warn(errors.As(err)); return }")
					*/
				}
				/*
					case *ast.StarExpr:
						// 指针类型
						switch tt := t.X.(type) {
						case *ast.SelectorExpr:
							// 其他包定义的指针类型
							g.p(name.Name, " := new(", tt.X, ".", tt.Sel, ")")
							// 导入其他包的定义
							g.pickupImports(fmt.Sprint(tt.X), imports)
						case *ast.Ident:
							// 当前包定义的指针类型
							g.p(name.Name, " := new(", tt.Name, ")")
						default:
							log.Printf("not supported type: %#v\n", tt)
							continue
						}
						g.p("if err := ctx.From", gt.Format, "(", toSnake(name.Name, gt.Snake), "); err != nil { log.Warn(errors.As(err)) return }")
				*/
				/*
					case *ast.SelectorExpr:
						// 其他包定义的非指针类型
						// 导入其他包的定义
						g.pickupImports(fmt.Sprint(t.X), imports)
						g.p("var ", name.Name, " ", t.X, ".", t.Sel)
						g.p("if err := ctx.From", gt.Format, "(&", toSnake(name.Name, gt.Snake), "); err != nil { log.Warn(errors.As(err)) return }")
				*/
			default:
				panic(fmt.Sprintf("not supported type: %#v\n : %s", t, name.Name))
				//log.Printf("not supported type: %#v\n", t)
				//continue
			}
			names = append(names, name.Name)
		}
	}

	if len(names) > 0 {
		g.p(fn.Name.Name, "(ctx, ", strings.Join(names, ","), ")")
	} else {
		g.p(fn.Name.Name, "(ctx)")
	}
	g.p("}")
	g.p()
}
