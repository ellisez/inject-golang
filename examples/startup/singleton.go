package startup

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/internal/vo"
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
// @import github.com/ellisez/inject-golang/examples/internal/vo
func NewConfig() *vo.Config {
	fmt.Println("call Config.NewConfig")
	return &vo.Config{
		Host: "127.0.0.1",
		Port: 3000,
	}
}

// PrepareDatabase example for explicit singleton
// @provide Database singleton
// @order "step 3: Setting Database"
// @import github.com/ellisez/inject-golang/examples/internal/vo
// @import github.com/ellisez/inject-golang/examples/web/handler
// @handler startup.DatabaseLoaded
func PrepareDatabase() *vo.Database {
	fmt.Println("call Database.PrepareDatabase")
	return &vo.Database{
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
// @provide ServerAlias _ vo.ServerInterface
// @order "step 4: Setting Server"
// @import github.com/ellisez/inject-golang/examples/internal/vo
// @injectParam config
// @injectParam database
// @injectParam argInt
// @handler ServerAliasLoaded
func PrepareServerAlias(config *vo.Config, database *vo.Database, argInt bool) *vo.Server {
	fmt.Println("call WebAppAlias.PrepareWebAppAlias")
	return &vo.Server{
		Config:   config,
		Database: database,
	}
}
