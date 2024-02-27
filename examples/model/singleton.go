package model

// Config
// @provide
type Config struct {
	Host string
	Port uint
}

// Database
// @provide _ singleton
// @import github.com/ellisez/inject-golang/examples/other
// @preConstruct model.PrepareDatabase
// @postConstruct other.DatabaseLoaded
type Database struct {
	Host     string
	Port     uint
	Schema   string
	UserName string
	Password string
}

func PrepareDatabase() *Database {
	return &Database{
		Host:     "127.0.0.1",
		Port:     3000,
		Schema:   "db",
		UserName: "admin",
		Password: "admin",
	}
}

// WebApp
// @provide WebCtxAlias
// @injectField Database
// @preConstruct model.PrepareWebCtxAlias
// @postConstruct WebCtxAliasLoaded
type WebApp struct {
	// @inject
	*Config
	*Database
	MiddleWares []*MiddleWare
	Routers     []*Router
}

func PrepareWebCtxAlias() *WebApp {
	// PreConstruct is usually used to set default values.
	// Unable to get Ctx, because it is not yet ready.
	// If you want to run with Ctx, you can use PostConstruct.
	/* error example
	var m = ctx.NewMiddleWare("/api", func() error {
		fmt.Printf("Call MiddleWare %s\n", "/api")
		return nil
	})
	return &WebApp{
		MiddleWares: []*MiddleWare{
			m,
		},
		Routers: []*Router{
			ctx.NewRouterAlias(m, "/login", func() error {
				fmt.Printf("Call Router %s", "/login")
				return nil
			}),
		},
	}*/
	return &WebApp{}
}
