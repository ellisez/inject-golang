package handler

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/model"
)

func DatabaseLoaded(database *model.Database) {
	fmt.Printf("call Database.postConstruct: %v\n", database)
}

func AfterRouterCreate(router *model.Router) {
	fmt.Printf("call Router.postConstruct: %v\n", router)
}
