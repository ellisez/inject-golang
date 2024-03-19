package provide

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/model"
)

// IsReady example for basic type
// @provide IsReady
// @order "step 1: Setting variable"
func IsReady() bool {
	fmt.Println("call IsReady.IsReady")
	return false
}

// NewConfig example for default singleton
// @provide Config
// @order "step 2: Load config"
// @import github.com/ellisez/inject-golang/examples/model
func NewConfig() *model.Config {
	fmt.Println("call Config.NewConfig")
	return &model.Config{
		Host: "127.0.0.1",
		Port: 3000,
	}
}

// PrepareDatabase example for explicit singleton
// @provide Database singleton
// @order "step 3: Setting Database"
// @import github.com/ellisez/inject-golang/examples/model
// @import github.com/ellisez/inject-golang/examples/handler
// @handler provide.DatabaseLoaded
func PrepareDatabase() *model.Database {
	fmt.Println("call Database.PrepareDatabase")
	return &model.Database{
		Host:     "127.0.0.1",
		Port:     3306,
		Schema:   "db",
		UserName: "admin",
		Password: "admin",
	}
}

func DatabaseLoaded() {
	fmt.Printf("call Database.DatabaseLoaded")
}

// PrepareServerAlias example for proxy handler
// @provide ServerAlias _ model.ServerInterface
// @order "step 4: Setting Server"
// @import github.com/ellisez/inject-golang/examples/model
// @injectParam config
// @injectParam database
// @injectParam argInt
// @handler ServerAliasLoaded
func PrepareServerAlias(config *model.Config, database *model.Database, argInt bool) *model.Server {
	fmt.Println("call WebAppAlias.PrepareWebAppAlias")
	return &model.Server{
		Config:   config,
		Database: database,
	}
}
