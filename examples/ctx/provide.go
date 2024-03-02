package ctx

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/init"
	"github.com/ellisez/inject-golang/examples/middleware"
	"github.com/ellisez/inject-golang/examples/model"
	"github.com/ellisez/inject-golang/examples/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"mime/multipart"
	"reflect"
	"strconv"
)

func Params(webCtx *fiber.Ctx, key string, defaultValue ...string) string {
	return webCtx.Params(key, defaultValue...)
}

func ParamsInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error) {
	return webCtx.ParamsInt(key, defaultValue...)
}

func ParamsBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error) {
	str := Params(webCtx, key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.ParseBool(str)
}

func ParamsFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error) {
	str := Params(webCtx, key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.ParseFloat(str, 64)
}

func ParamsParser(webCtx *fiber.Ctx, out any) error {
	return webCtx.ParamsParser(out)
}

func Query(webCtx *fiber.Ctx, key string, defaultValue ...string) string {
	return webCtx.Query(key, defaultValue...)
}

func QueryInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error) {
	str := webCtx.Query(key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.Atoi(str)
}

func QueryBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error) {
	str := webCtx.Query(key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.ParseBool(str)
}

func QueryFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error) {
	str := webCtx.Query(key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.ParseFloat(str, 64)
}

func QueryParser(webCtx *fiber.Ctx, out any) error {
	return webCtx.QueryParser(out)
}

func Header(webCtx *fiber.Ctx, key string, defaultValue ...string) string {
	return webCtx.GetRespHeader(key, defaultValue...)
}

func HeaderInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error) {
	str := Header(webCtx, key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.Atoi(str)
}

func HeaderBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error) {
	str := Header(webCtx, key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.ParseBool(str)
}

func HeaderFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error) {
	str := Header(webCtx, key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.ParseFloat(str, 64)
}

func HeaderParser(webCtx *fiber.Ctx, out any) error {
	return webCtx.ReqHeaderParser(out)
}

func FormString(webCtx *fiber.Ctx, key string, defaultValue ...string) string {
	return webCtx.FormValue(key, defaultValue...)
}
func FormFile(webCtx *fiber.Ctx, key string) (*multipart.FileHeader, error) {
	return webCtx.FormFile(key)
}

func FormInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error) {
	str := FormString(webCtx, key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.Atoi(str)
}

func FormBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error) {
	str := FormString(webCtx, key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.ParseBool(str)
}

func FormFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error) {
	str := FormString(webCtx, key)
	if str == "" && defaultValue != nil {
		return defaultValue[0], nil
	}
	return strconv.ParseFloat(str, 64)
}

func FormParser(webCtx *fiber.Ctx, out any) error {
	elem := reflect.ValueOf(out).Elem()
	form, err := webCtx.MultipartForm()
	if err != nil {
		return err
	}
	for key, strArr := range form.Value {
		for _, value := range strArr {
			field := elem.FieldByName(key)
			if field.IsValid() && field.CanSet() {
				switch field.Kind() {
				case reflect.String:
					field.SetString(value)
				case reflect.Int:
					intValue, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					field.SetInt(int64(intValue))
				case reflect.Bool:
					boolValue, err := strconv.ParseBool(value)
					if err != nil {
						return err
					}
					field.SetBool(boolValue)
				case reflect.Float64:
					floatValue, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return err
					}
					field.SetFloat(floatValue)
				default:
					return fmt.Errorf("unsupported type %T", value)
				}
			} else if !field.IsValid() {
				// 如果没有找到对应的字段则返回错误
			}
		}
	}

	for key, fileArr := range form.File {
		for _, file := range fileArr {
			field := elem.FieldByName(key)
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(file))
			} else if !field.IsValid() {
				// 如果没有找到对应的字段则返回错误
			}
		}
	}
	return nil
}

func Body(webCtx *fiber.Ctx) []byte {
	return webCtx.Body()
}
func BodyString(webCtx *fiber.Ctx) string {
	return string(Body(webCtx))
}
func BodyParser(webCtx *fiber.Ctx, out any) error {
	return webCtx.BodyParser(out)
}

func (ctx *Ctx) WebAppStartup() error {

	webApp := fiber.New()
	ctx.WebApp = webApp

	/// middlewares
	/// routers
	/// groups {
	///     - /// middlewares
	///	    - /// router
	///      - /// group
	/// }
	webApp.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))
	webApp.Static("/", "/", fiber.Static{
		Compress: true,
		Download: true,
		Browse:   true,
		Index:    "",
	})

	webApp.Group("/api", ctx.CorsMiddleware)

	webApp.Post("/api/login", ctx.LoginController)

	host, port, err := init.ConfigureWebApp(webApp, ctx.Config)
	if err != nil {
		return err
	}
	return webApp.Listen(fmt.Sprintf("%s:%d", host, port))

}

func (ctx *Ctx) CorsMiddleware(webCtx *fiber.Ctx) error {
	config := &model.Config{}
	err := BodyParser(webCtx, config)
	if err != nil {
		return err
	}

	header := Header(webCtx, "header")

	paramsInt, err := ParamsInt(webCtx, "paramsInt")
	if err != nil {
		return err
	}

	queryBool, err := QueryBool(webCtx, "queryBool")
	if err != nil {
		return err
	}

	formFloat, err := FormFloat(webCtx, "formFloat")
	if err != nil {
		return err
	}
	return middleware.CorsMiddleware(webCtx, config, header, paramsInt, queryBool, formFloat)
}

func (ctx *Ctx) LoginController(webCtx *fiber.Ctx) (err error) {
	/// body -> BodyParser -> []
	config := &model.Config{}

	bodyJson := make(map[string]any)
	err = webCtx.BodyParser(bodyJson)
	if err != nil {
		return err
	}
	username := bodyJson["username"].(string)
	password := bodyJson["password"].(string)

	/// Path
	Params(webCtx, "path")
	ParamsBool(webCtx, "")
	ParamsFloat(webCtx, "")
	ParamsInt(webCtx, "path")
	ParamsParser(webCtx, config)

	/// query
	Query(webCtx, "path")
	QueryBool(webCtx, "")
	QueryFloat(webCtx, "")
	QueryInt(webCtx, "path")
	QueryParser(webCtx, config)

	/// header
	Header(webCtx, "header")
	HeaderBool(webCtx, "")
	HeaderFloat(webCtx, "")
	HeaderInt(webCtx, "path")
	HeaderParser(webCtx, config)

	/// formData
	FormString(webCtx, "formData")
	FormBool(webCtx, "formData")
	FormFloat(webCtx, "formData")
	FormInt(webCtx, "formData")
	FormFile(webCtx, "formData")
	FormParser(webCtx, config)

	/// multipart
	form, err := webCtx.MultipartForm()
	if err != nil {
		return err
	}

	var multipartField string
	multipartFields := form.Value["multipartField"]
	if multipartFields != nil && len(multipartFields) == 1 {
		multipartField = multipartFields[0]
	}

	var multipartFile *multipart.FileHeader
	multipartFiles := form.File["multipartFile"]
	if multipartFiles != nil && len(multipartFiles) == 1 {
		multipartFile = multipartFiles[0]
	}

	fmt.Println(multipartField, multipartFile)
	return router.LoginController(username, password)
}
