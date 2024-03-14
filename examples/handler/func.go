package handler

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/ctx"
	"github.com/ellisez/inject-golang/examples/model"
)

// ServerAliasLoaded example for injection proxy
// @proxy
// @injectParam database Database
// @injectCtx appCtx
// @injectParam isReady _ &
// @injectParam event
// @injectParam listener
func ServerAliasLoaded(appCtx ctx.Ctx, server *model.Server, database *model.Database, isReady *bool, event *model.Event, listener *model.Listener) {
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

func AfterRouterCreate(router *model.Listener) {
	fmt.Printf("call Listener.postConstruct: %v\n", router)
}
