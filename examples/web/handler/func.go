package handler

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/ctx"
	"github.com/ellisez/inject-golang/examples/internal/vo"
)

// ServerAliasLoaded example for injection proxy
// @proxy
// @import "github.com/ellisez/inject-golang/examples/model"
// @injectParam database Database
// @injectCtx appCtx
// @injectParam server ServerAlias cast
// @injectParam isReady _ &
// @injectParam event
// @injectParam listener
func ServerAliasLoaded(appCtx ctx.Ctx, server *vo.Server, database *vo.Database, isReady *bool, event *vo.Event, listener *vo.Listener) {
	fmt.Printf("call proxy.WebAppAliasLoaded: %v, %v, %v\n", server, database, isReady)
	server.Startup()
	*isReady = true
	appCtx.TestServer(server)
	// custom
	server.AddListener("register", func(data map[string]any) {
		fmt.Printf("call Event: '%s', Listener: %v\n", "register", data)
	})
	server.AddListener("login", func(data map[string]any) {
		fmt.Printf("call Event: '%s', Listener: %v\n", "register", data)
	})
}

func AfterRouterCreate() {
	fmt.Printf("call Listener.postConstruct:\n")
}
