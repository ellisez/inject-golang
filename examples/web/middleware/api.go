package middleware

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/ctx"
	"github.com/ellisez/inject-golang/examples/internal/vo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CorsMiddleware
// @middleware /api
// @import github.com/ellisez/inject-golang/examples/model
// @injectWebCtx c
// @injectCtx appCtx
// @param body body
// @param header header
// @param paramsInt path
// @param queryBool query
// @param formFloat formData
func CorsMiddleware(appCtx ctx.Ctx, c *fiber.Ctx,
	body *vo.Config,
	header string,
	paramsInt int,
	queryBool bool,
	formFloat float64,
) error {
	fmt.Printf("call CorsMiddleware: %v, %v, %v, %s, %d, %t, %f\n", appCtx, c, body, header, paramsInt, queryBool, formFloat)
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	})(c)
}
