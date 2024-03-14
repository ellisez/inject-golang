package provide

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/model"
)

// NewEvent example for multiple default
// @provide Event multiple
// @import github.com/ellisez/inject-golang/examples/model
// @injectParam Database
// @injectParam config
func NewEvent(Database *model.Database, config *model.Config) *model.Event {
	fmt.Println("call Event.NewEvent")
	return &model.Event{
		Config:   config,
		Database: Database,
	}
}

// NewListener example for multiple injection
// @provide Listener multiple
// @import github.com/ellisez/inject-golang/examples/model
// @import github.com/ellisez/inject-golang/examples/handler
// @handler handler.AfterRouterCreate
func NewListener() *model.Listener {
	fmt.Println("call Listener.NewListener")
	return &model.Listener{
		Func: func(_ map[string]any) {
		},
	}
}
