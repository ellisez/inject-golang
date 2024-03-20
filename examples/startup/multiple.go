package startup

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/internal/vo"
)

// NewEvent example for multiple default
// @provide Event multiple
// @import github.com/ellisez/inject-golang/examples/model
// @injectParam Database
// @injectParam config
func NewEvent(Database *vo.Database, config *vo.Config) *vo.Event {
	fmt.Println("call Event.NewEvent")
	return &vo.Event{
		Config:   config,
		Database: Database,
	}
}

// NewListener example for multiple injection
// @provide Listener multiple
// @import github.com/ellisez/inject-golang/examples/model
// @import github.com/ellisez/inject-golang/examples/web/handler
// @handler handler.AfterRouterCreate
func NewListener() *vo.Listener {
	fmt.Println("call Listener.NewListener")
	return &vo.Listener{
		Func: func(_ map[string]any) {
		},
	}
}
