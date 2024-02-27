package model

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
// @preConstruct model.NewRouterAlias
type Router struct {
	MiddleWare *MiddleWare
	Path       string
	Handle     func() error
}

func NewRouterAlias() *Router {
	return &Router{}
}
