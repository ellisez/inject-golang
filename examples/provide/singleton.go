package provide

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/model"
)

// ProvideIsReady example for basic type
// @provide IsReady
// @order "step 1: Setting variable"
func ProvideIsReady() bool {
	fmt.Println("call IsReady.ProvideIsReady")
	return false
}

// NewConfig example for default singleton
// @provide Config
// @order "step 2: Load config"
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

func DatabaseLoaded(database *model.Database) {
	fmt.Printf("call Database.DatabaseLoaded: %v\n", database)
}

// PrepareWebAppAlias example for proxy handler
// @provide WebAppAlias
// @order "step 4: Setting WebApp"
// @handler WebAppAliasLoaded
func PrepareWebAppAlias() *model.WebApp {
	fmt.Println("call WebAppAlias.PrepareWebAppAlias")
	return &model.WebApp{}
}
