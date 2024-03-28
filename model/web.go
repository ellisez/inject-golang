package model

import (
	"go/ast"
)

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
	Methods []string // 请求方式get|head|post|put|patch|delete|connect|options|trace
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
	keys        []*Key
	resources   map[string]*WebResource // 静态资源
	middlewares map[string]*Middleware  // 组内中间件
	routers     map[string]*Router      // 组内路由
}

func NewWebApplication() *WebApplication {
	webInstance := &WebApplication{}
	webInstance.resources = map[string]*WebResource{}
	webInstance.middlewares = map[string]*Middleware{}
	webInstance.routers = map[string]*Router{}
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
			Alias: "",
			Path:  "github.com/gofiber/fiber/v2",
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

func (w *WebApplication) IndexOf(index int) *Key {
	return w.keys[index]
}
func (w *WebApplication) KeyOf(key string) *Key {
	for _, k := range w.keys {
		if k.Instance == key {
			return k
		}
	}
	return nil
}

//

func (w *WebApplication) AddMiddleware(middleware *Middleware) {
	k := &Key{
		Instance: middleware.Instance,
		Type:     "middleware",
		Order:    middleware.Order,
	}
	w.keys = append(w.keys, k)
	w.middlewares[middleware.Instance] = middleware
}

func (w *WebApplication) AddRouter(router *Router) {
	k := &Key{
		Instance: router.Instance,
		Type:     "router",
		Order:    router.Order,
	}
	w.keys = append(w.keys, k)
	w.routers[router.Instance] = router
}

func (w *WebApplication) AddResource(resource *WebResource) {
	k := &Key{
		Instance: resource.Path,
		Type:     "resource",
		Order:    "",
	}
	w.keys = append(w.keys, k)
	w.resources[resource.Path] = resource
}

func (w *WebApplication) DelMiddleware(key string) {
	var hitKey *Key
	for i, k := range w.keys {
		if k.Instance == key && k.Type == "middleware" {
			w.keys = append(w.keys[:i], w.keys[i+1:]...)
			hitKey = k
			break
		}
	}
	if hitKey != nil {
		delete(w.middlewares, key)
	}
}

func (w *WebApplication) DelRouter(key string) {
	var hitKey *Key
	for i, k := range w.keys {
		if k.Instance == key && k.Type == "router" {
			w.keys = append(w.keys[:i], w.keys[i+1:]...)
			hitKey = k
			break
		}
	}
	if hitKey != nil {
		delete(w.routers, key)
	}
}

func (w *WebApplication) DelResource(key string) {
	var hitKey *Key
	for i, k := range w.keys {
		if k.Instance == key && k.Type == "resource" {
			w.keys = append(w.keys[:i], w.keys[i+1:]...)
			hitKey = k
			break
		}
	}
	if hitKey != nil {
		delete(w.resources, key)
	}
}
func (w *WebApplication) ContainsMiddleware(key string) bool {
	return w.GetMiddleware(key) != nil
}

func (w *WebApplication) ContainsRouter(key string) bool {
	return w.GetRouter(key) != nil
}

func (w *WebApplication) ContainsResource(key string) bool {
	return w.GetResource(key) != nil
}

func (w *WebApplication) GetMiddleware(instance string) *Middleware {
	return w.middlewares[instance]
}

func (w *WebApplication) GetRouter(instance string) *Router {
	return w.routers[instance]
}
func (w *WebApplication) GetResource(instance string) *WebResource {
	return w.resources[instance]
}

func (w *WebApplication) ReplaceMiddleware(middleware *Middleware) bool {
	m := w.middlewares[middleware.Instance]
	if m == nil {
		return false
	}
	key := w.KeyOf(middleware.Instance)
	key.Order = middleware.Order
	w.middlewares[middleware.Instance] = middleware
	return true
}

func (w *WebApplication) ReplaceRouter(router *Router) bool {
	m := w.routers[router.Instance]
	if m == nil {
		return false
	}
	key := w.KeyOf(router.Instance)
	key.Order = router.Order
	w.routers[router.Instance] = router
	return true
}
func (w *WebApplication) ReplaceResource(resource *WebResource) bool {
	m := w.resources[resource.Path]
	if m == nil {
		return false
	}
	w.resources[resource.Path] = resource
	return true
}

func (w *WebApplication) Merge(other *WebApplication) {
	for _, resource := range other.resources {
		_, has := w.resources[resource.Path]
		if has {
			w.ReplaceResource(resource)
		} else {
			w.AddResource(resource)
		}
	}
	for _, middleware := range other.middlewares {
		_, has := w.middlewares[middleware.Instance]
		if has {
			w.ReplaceMiddleware(middleware)
		} else {
			w.AddMiddleware(middleware)
		}
	}
	for _, router := range other.routers {
		_, has := w.middlewares[router.Instance]
		if has {
			w.ReplaceRouter(router)
		} else {
			w.AddRouter(router)
		}
	}
}

func (w *WebApplication) MiddlewareLen() int {
	return len(w.middlewares)
}
func (w *WebApplication) RouterLen() int {
	return len(w.routers)
}
func (w *WebApplication) ResourceLen() int {
	return len(w.resources)
}

// for sort

func (w *WebApplication) Len() int {
	return len(w.keys)
}

func (w *WebApplication) Swap(x int, y int) {
	old := w.keys[x]
	w.keys[x] = w.keys[y]
	w.keys[y] = old
}
func (w *WebApplication) Less(x int, y int) bool {
	a := w.keys[x]
	b := w.keys[y]
	orderA := a.Order
	orderB := b.Order
	if orderA != "" && orderB == "" {
		return true
	}
	if orderA == "" && orderB != "" {
		return false
	}
	return orderA < orderB
}
