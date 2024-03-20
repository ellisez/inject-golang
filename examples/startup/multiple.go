package startup

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/internal"
)

// NewEvent example for multiple default
// @provide Event multiple
// @import github.com/ellisez/inject-golang/examples/model
// @injectParam Database
// @injectParam config
func NewEvent(Database *internal.Database, config *internal.Config) *internal.Event {
	fmt.Println("call Event.NewEvent")
	return &internal.Event{
		Config:   config,
		Database: Database,
	}
}

// NewListener example for multiple injection
// @provide Listener multiple
// @import github.com/ellisez/inject-golang/examples/model
// @import github.com/ellisez/inject-golang/examples/web/handler
// @handler handler.AfterRouterCreate
func NewListener() *internal.Listener {
	fmt.Println("call Listener.NewListener")
	return &internal.Listener{
		Func: func(_ map[string]any) {
		},
	}
}
