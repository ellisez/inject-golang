package middleware

import (
	"github.com/ellisez/inject-golang/examples/ctx"
	"github.com/ellisez/inject-golang/examples/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CorsMiddleware
// @middleware /api
// @injectWebCtx c
// @injectCtx appCtx
// @param body body
// @param header header
// @param paramsInt path
// @param queryBool query
// @param formFloat formData
func CorsMiddleware(appCtx *ctx.Ctx, c *fiber.Ctx,
	body *model.Config,
	header string,
	paramsInt int,
	queryBool bool,
	formFloat float64) error {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	})(c)
}
