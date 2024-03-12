package provide

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/model"
)

// NewMiddleWare example for multiple default
// @provide MiddleWare multiple
// @injectParam Database
func NewMiddleWare(Database *model.Database) *model.MiddleWare {
	return &model.MiddleWare{
		Database: Database,
	}
}

// NewRouterAlias example for multiple injection
// @provide RouterAlias multiple
// @import github.com/ellisez/inject-golang/examples/handler
// @handler handler.AfterRouterCreate
func NewRouterAlias() *model.Router {
	fmt.Println("call Router.NewRouterAlias")
	return &model.Router{
		Path: "/login",
		Handle: func() error {
			fmt.Println("call Router.Handle")
			return nil
		},
	}
}
