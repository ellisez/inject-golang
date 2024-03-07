package model

import "fmt"

// Config
// @provide
type Config struct {
	Host string
	Port uint
}

// Database
// @provide _ singleton
// @import github.com/ellisez/inject-golang/examples/handler
// @preConstruct model.PrepareDatabase
// @postConstruct handler.DatabaseLoaded
type Database struct {
	Host     string
	Port     uint
	Schema   string
	UserName string
	Password string
}

func PrepareDatabase() *Database {
	fmt.Println("call Database.preConstruct")
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
// @import github.com/ellisez/inject-golang/examples/handler
// @injectField Database
// @preConstruct handler.PrepareWebCtxAlias
// @postConstruct WebCtxAliasLoaded
type WebApp struct {
	// @inject
	*Config
	*Database
	MiddleWares []*MiddleWare
	Routers     []*Router
}

// TestLogin
// @proxy
// @injectParam database Database
func (webApp *WebApp) TestLogin(database *Database) {
	fmt.Printf("call TestLogin: %v, %v\n", webApp, database)
}
