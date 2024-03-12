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
// @order first
// @import github.com/ellisez/inject-golang/examples/handler
// @injectField database Database _ _
// @preConstruct handler.PrepareWebCtxAlias
// @postConstruct WebCtxAliasLoaded
type WebApp struct {
	// @inject Config Config SetConfig1
	config      *Config
	database    *Database
	MiddleWares []*MiddleWare
	Routers     []*Router
}

func (w *WebApp) SetDatabase(database *Database) {
	fmt.Println("call instance.database setter")
	w.database = database
}

func (w *WebApp) Database() *Database {
	fmt.Println("call instance.database getter")
	return w.database
}

func (w *WebApp) Config() *Config {
	fmt.Println("call instance.config getter")
	return w.config
}

func (w *WebApp) SetConfig1(config *Config) {
	fmt.Println("call instance.config setter")
	w.config = config
}

// TestLogin
// @proxy
// @injectParam database Database
func (w *WebApp) TestLogin(database *Database) {
	fmt.Printf("call TestLogin: %v, %v\n", w, database)
}
