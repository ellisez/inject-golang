package handler

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/ctx"
	"github.com/ellisez/inject-golang/examples/model"
)

// WebCtxAliasLoaded
// @proxy
// @import github.com/ellisez/inject-golang/examples/model
// @injectParam database Database
// @injectParam router RouterAlias
// @injectCtx appCtx
func WebCtxAliasLoaded(appCtx ctx.Ctx, webApp *model.WebApp, database *model.Database, router *model.Router) {
	fmt.Printf("call WebApp.postConstruct: %v, %v, %v\n", webApp, database, router)
	appCtx.TestLogin(webApp)
}
