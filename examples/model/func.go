package model

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/ctx"
)

// WebCtxAliasLoaded
// @proxy
// @injectParam database Database
// @injectParam ctx
func WebCtxAliasLoaded(ctx *ctx.Ctx, webApp *WebApp, database *Database) {
	fmt.Printf("WebCtxAliasLoaded: %v\n%v\n", webApp, database)
	ctx.TestLogin(webApp)
}

// TestLogin
// @proxy
// @injectParam database Database
func (webApp *WebApp) TestLogin(database *Database) {
	fmt.Printf("WebCtxAliasLoaded: %v\n%v\n", webApp, database)
}