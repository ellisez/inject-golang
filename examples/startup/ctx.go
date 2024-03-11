package startup

import "github.com/ellisez/inject-golang/examples/ctx"

// CtxConfigure
// @ctxProvide
// @proxy
// @injectCtx appCtx
// @value Welcome string "hello world!"
func CtxConfigure(appCtx ctx.Ctx) {
	appCtx.SetWelcome("hello inject-golang")
}
