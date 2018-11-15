package ginIOC_param

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

//pkg where handlers to generate
var pkg string 
//Request Body format, default json.
var content = "json" 
//param name ,is snake or camelcase, default camelcase.
var snake = true 
 //pkg of frame
var handlerPkgName []string=[]string{"github.com/gin-gonic/gin","tool/gintool"}
//pkg of log
var logPkgName string="github.com/qiniu/log"
//pkg of error
var errorsPkgName string="tool/errors" 
//pkg of this
var thisPkgName string=	"tool/ginIOC_param"

type generator struct {
	*bytes.Buffer
	pkg          *Package
	copyBuffer   *bytes.Buffer
	headerBuffer *bytes.Buffer
	initBuffer   *bytes.Buffer
	imports      map[string]string
}

type Package struct {
	name  string
	files []*ast.File
	fset  *token.FileSet
}

//handlerPkg:pkg where handler handlers to generate 
func IOCParam(handlerPkg string) {
	pkg=handlerPkg
	var g = &generator{
		Buffer:       new(bytes.Buffer),
		copyBuffer:   new(bytes.Buffer),
		headerBuffer: new(bytes.Buffer),
		initBuffer:   new(bytes.Buffer),
		imports:      make(map[string]string),
	}
	g.parsePackageDir(pkg)
	g.do(filepath.Join(pkg, "gen_handler.go"), g.pluginHandler)
	// Run generate.
}

func (g *generator) do(file string, f func()) {
	g.Reset()
	g.copyBuffer.Reset()
	g.headerBuffer.Reset()
	g.initBuffer.Reset()
	g.imports = make(map[string]string)

	f()

	// Print the header.
	g.pHeader(fmt.Sprintf(`// Code generated by "gen -pkg %s"; DO NOT EDIT.`, pkg))
	g.pHeader()
	g.pHeader("package ", g.pkg.name)
	g.pHeader()
	if len(g.imports) > 0 {
		g.pHeader("import (")
		for pkg, alias := range g.imports {
			if alias == "" {
				g.pHeader(`"`, pkg, `"`)
			} else if !strings.Contains(pkg, `"`) {
				g.pHeader(alias, ` "`, pkg, `"`)
			} else {
				g.pHeader(alias, " ", pkg)
			}
		}
		g.pHeader(")")
		g.pHeader()
	}

	if bytes := g.initBuffer.Bytes(); len(bytes) > 0 {
		g.p()
		g.p("func init() {")
		g.Write(bytes)
		g.p("}")
	}

	// Format the output.
	src := g.format()
	if len(src) == 0 {
		return
	}

	if err := ioutil.WriteFile(file, src, 0644); err != nil {
		panic(fmt.Sprintf("writing output: %s\n", err))
	}
}

func (g *generator) importPkg(pkg, alias string) {
	g.imports[pkg] = alias
}

func (g *generator) pBuffer(buf *bytes.Buffer, args ...interface{}) {
	for _, v := range args {
		switch s := v.(type) {
		case string:
			buf.WriteString(s)
		case *string:
			buf.WriteString(*s)
		case bool:
			buf.WriteString(fmt.Sprintf("%t", s))
		case *bool:
			buf.WriteString(fmt.Sprintf("%t", *s))
		case int:
			buf.WriteString(fmt.Sprintf("%d", s))
		case *int32:
			buf.WriteString(fmt.Sprintf("%d", *s))
		case *int64:
			buf.WriteString(fmt.Sprintf("%d", *s))
		case float64:
			buf.WriteString(fmt.Sprintf("%g", s))
		case *float64:
			buf.WriteString(fmt.Sprintf("%g", *s))
		default:
			panic(fmt.Sprintf("warning: unknown type in printer: %T\n", v))
		}
	}
	buf.WriteByte('\n')
}

func (g *generator) p(args ...interface{}) {
	g.pBuffer(g.Buffer, args...)
}

func (g *generator) pHeader(args ...interface{}) {
	g.pBuffer(g.headerBuffer, args...)
}

func (g *generator) pInit(args ...interface{}) {
	g.pBuffer(g.initBuffer, args...)
}

func (g *generator) parsePackage(directory string, names []string) {
	g.pkg = new(Package)
	g.pkg.fset = token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") || strings.HasPrefix(name, "gen_") {
			continue
		}
		parsedFile, err := parser.ParseFile(g.pkg.fset, name, nil, parser.ParseComments)
		if err != nil {
			panic(fmt.Sprintf("parsing package: %s: %s\ns", name, err))
		}
		g.pkg.files = append(g.pkg.files, parsedFile)
	}
	if len(g.pkg.files) == 0 {
		panic(fmt.Sprintf("%s: no buildable Go files", directory))
	}
	g.pkg.name = g.pkg.files[0].Name.Name
}

func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}

func (g *generator) parsePackageDir(dir string) {
	pkg, err := build.Default.ImportDir(dir, 0)
	if err != nil {
		panic(fmt.Sprintf("cannot process directory %s: %s\n", dir, err))
	}
	g.parsePackage(dir, prefixDirectory(dir, pkg.GoFiles))
}

type genTag struct {
	Decl   ast.Decl
	Gs     []*ast.CommentGroup
	Tag    string
	Format string
	Snake  bool
}

type genTags []*genTag

func (gts genTags) Len() int {
	return len(gts)
}

func (gts genTags) Less(i, j int) bool {
	return gts[i].Decl.Pos() < gts[j].Decl.Pos()
}

func (gts genTags) Swap(i, j int) {
	gts[i], gts[j] = gts[j], gts[i]
}

func (g *generator) pluginHandler() {
	for _, f := range g.pkg.files {
		cmap := ast.NewCommentMap(g.pkg.fset, f, f.Comments)
		handlers := make([]*genTag, 0)
		for key, val := range cmap {
			if fn, ok := key.(*ast.FuncDecl); ok {
				for _, commentGroup := range val {
					for _, comment := range commentGroup.List {
						if strings.HasPrefix(comment.Text, "// @handler") {
							gt := &genTag{
								Decl:   fn,
								Gs:     val,
								Tag:    comment.Text,
								Format: strings.Title(strings.ToLower(content)),
								Snake:  snake,
							}
							handlers = append(handlers, gt)
							break
						}
					}
				}
			}
		}
		var imports = make(map[string]string)
		sort.Sort(genTags(handlers))
		if len(handlers) > 0 {
			for _,v:=range handlerPkgName{
				g.importPkg(v, "")
			}
			for _, imp := range f.Imports {
				var alias string
				if imp.Name != nil {
					alias = imp.Name.Name
				}
				if imp.Path != nil {
					imports[imp.Path.Value] = alias
				}
			}
		}
		for _, gt := range handlers {
			g.generateHandler(gt, imports)
		}
	}
}

func (g *generator) format() []byte {
	origin := g.headerBuffer.Bytes()
	body := g.Bytes()
	if len(body) == 0 {
		return nil
	}
	origin = append(origin, body...)
	src, err := format.Source(origin)
	if err != nil {
		//panic(err)
		panic(fmt.Sprintf("warning: internal error: invalid Go generated: %s\n", err))
		return origin
	}
	return src
}