package model

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/ctx"
)

// WebCtxAliasLoaded
// @proxy
// @injectParam database Database
// @injectParam router RouterAlias
// @injectCtx appCtx
func WebCtxAliasLoaded(appCtx *ctx.Ctx, webApp *WebApp, database *Database, router *Router) {
	fmt.Printf("call WebApp.postConstruct: %v\n%v\n%v\n", webApp, database, router)
	appCtx.TestLogin(webApp)
}

// TestLogin
// @proxy
// @injectParam database Database
func (webApp *WebApp) TestLogin(database *Database) {
	fmt.Printf("call TestLogin: %v\n%v\n", webApp, database)
}
