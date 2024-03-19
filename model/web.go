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

type WebApplication struct {
	Resources   map[string]*WebResource // 静态资源
	Middlewares map[string]*Middleware  // 组内中间件
	Routers     map[string]*Router      // 组内路由
}

func NewWebApplication() *WebApplication {
	webInstance := &WebApplication{}
	webInstance.Resources = map[string]*WebResource{}
	webInstance.Middlewares = map[string]*Middleware{}
	webInstance.Routers = map[string]*Router{}
	return webInstance
}

func NewWebProvide() *Provide {
	provide := NewProvide()
	provide.Mode = "singleton"
	provide.Instance = "WebApp"
	provide.Type = &ast.StarExpr{
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("fiber"),
			Sel: ast.NewIdent("App"),
		}}
	provide.Imports = []*Import{
		{
			Name: "",
			Path: "github.com/gofiber/fiber/v2",
		},
	}
	provide.Constructor = &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("fiber"),
			Sel: ast.NewIdent("New"),
		},
	}
	return provide
}
