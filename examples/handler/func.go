package handler

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/ctx"
	"github.com/ellisez/inject-golang/examples/model"
)

// WebAppAliasLoaded example for injection proxy
// @proxy
// @injectParam database Database
// @injectCtx appCtx
// @injectParam webApp WebAppAlias
// @injectParam isReady IsReady &
// @injectParam middleware MiddleWare
// @injectParam router RouterAlias
func WebAppAliasLoaded(appCtx ctx.Ctx, webApp *model.WebApp, database *model.Database, isReady *bool, middleware *model.MiddleWare, router *model.Router) {
	fmt.Printf("call proxy.WebAppAliasLoaded: %v, %v, %v\n", webApp, database, isReady)
	webApp.Database = database
	webApp.Config = appCtx.Config()
	webApp.MiddleWares = append(webApp.MiddleWares, middleware, appCtx.NewMiddleWare())
	antherRouter := appCtx.NewRouterAlias()
	antherRouter.Path = "/logout"
	webApp.Routers = append(webApp.Routers, router, antherRouter)
	appCtx.TestLogin(webApp)
	*isReady = true
}

func AfterRouterCreate(router *model.Router) {
	fmt.Printf("call Router.postConstruct: %v\n", router)
}
