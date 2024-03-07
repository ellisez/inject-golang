package model

import "fmt"

// MiddleWare
// @provide _ multiple
// @injectField Database
type MiddleWare struct {
	// @inject
	Config *Config
	*Database
	Path   string
	Handle func() error
}

// Router
// @provide RouterAlias multiple
// @import github.com/ellisez/inject-golang/examples/handler
// @preConstruct model.NewRouterAlias
// @postConstruct handler.AfterRouterCreate
type Router struct {
	MiddleWare *MiddleWare
	Path       string
	Handle     func() error
}

func NewRouterAlias() *Router {
	fmt.Println("call Router.preConstruct")
	return &Router{}
}
