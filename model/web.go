package model

import "go/ast"

type WebParam struct {
	QueryParams  []*Field // query取值
	PathParams   []*Field // path取值
	HeaderParams []*Field // header取值
	BodyParam    *Field   // body取值
	FormParams   []*Field // formData取值
}

func NewWebParam() *WebParam {
	return &WebParam{}
}

type Middleware struct {
	*Proxy

	*WebParam
	WebApp string // 所属WebApp，默认WebApp
	Path   string // 路径
}

func NewMiddleware() *Middleware {
	return &Middleware{Proxy: NewProxy(), WebParam: NewWebParam()}
}

type Router struct {
	*Middleware
	Methods []string // 请求方式get|post|put|patch
}

func NewRouter() *Router {
	return &Router{Middleware: NewMiddleware()}
}

type WebResource struct {
	Path     string
	Dirname  string
	Features []string
	Index    string
	MaxAge   int
	Comment  string
}

type WebInstance struct {
	*Provide
	Resources   []*WebResource // 静态资源
	Middlewares []*Middleware  // 组内中间件
	Routers     []*Router      // 组内路由
}

func NewWebInstance() *WebInstance {
	webInstance := &WebInstance{Provide: NewProvide()}
	webInstance.Mode = "singleton"
	webInstance.Package = "fiber"
	webInstance.FuncName = "New"
	webInstance.Instance = "WebApp"
	webInstance.Type = &ast.StarExpr{
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("fiber"),
			Sel: ast.NewIdent("App"),
		}}
	webInstance.Imports = []*Import{
		{
			Name: "",
			Path: "github.com/gofiber/fiber/v2",
		},
	}
	return webInstance
}
