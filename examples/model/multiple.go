package model

import "fmt"

// MiddleWare
// @provide _ multiple
// @injectField Database
type MiddleWare struct {
	// @inject Config
	config *Config
	*Database
	Path   string
	Handle func() error
}

func (m *MiddleWare) Config() *Config {
	fmt.Println("call MiddleWare.config getter")
	return m.config
}

func (m *MiddleWare) SetConfig(config *Config) {
	fmt.Println("call MiddleWare.config setter")
	m.config = config
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
