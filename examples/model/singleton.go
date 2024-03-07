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
	fmt.Println("call WebApp.preConstruct")
	return &WebApp{}
}
