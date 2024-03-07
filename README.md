# inject-golang
Provide container of DI(Dependency Injection) for golang.

Language: English [中文](README_cn.md) 

### 1. Install and Run

```shell
go install github.com/ellisez/inject-golang
```
### 1.1. use go generate
> By default, scan the current module and generate all annotations.
```go
package main

import (
	"github.com/ellisez/inject-golang/examples/ctx"
)

//go:generate inject-golang
func main() {
	ctx.New()
}
```
> Use '- m' to specify that only partial annotations are generated, support setting multiple, separated by commas, and optional values such as 'singleton','multiple','fun', and 'web'

```go
//go:generate inject-golang -m singleton,multiple
func main() {
    ctx.New()
}
```

> Can specify scanning modules: supports multiple, defaults to the current module, external modules must be imported.
>
> "." indicates current package, The system already supports go.work.

```go
//go:generate inject-golang -m singleton,web github.com/ellisez/inject-golang/examples-work .
func main() {
    ctx.New()
}
```


To learn more about commands, please run `inject-golang -h`.

### 1.2. Run
```shell
go generate -run inject-golang
```

### 1.3. Clean Code
```shell
inject-glang --clean
```

## 2. Annotates
### 2.1. Struct Annotate
```
// @provide <Instance, default structName> <singleton default|multiple>
// @import *<Path, required> <Alias>
// @injectField *<FieldName, required> <Instance，default structName>
// @preConstruct *<before create call func, required>
// @postConstruct *<after created call func, required>
```

> The system will load the container package and the package where the annotation is located by default, but if there are additional packages, they need to be declared with `@import`;
>
> `@preConstruct` is a function without parameters but returns the type of structure, mainly used to create the initial value of the given structure before creation, and does not support injecting dependencies;
>
> `@postConstruct` is a function with a unique parameter of the structure type, mainly used to retrieve instances of structures that have already been created;
>
> `@preConstruct` and `@postConstruct` can carry package names, but they need to be used in conjunction with `@import`, such as "*model.Database";
>
> We generally do not recommend using the original function directly, but rather its proxy function, which can help us expand the parameters of other dependency injections;
> Proxy function usage, please refer to [Func Annotate (use for all func)](#23-func-annotate-use-for-all-func)
>
>

### 2.2. Field Annotate in struct
```
// @inject <Instance，default fieldName>
```
### 2.3. Func Annotate (use for all func)
```
// @proxy <Instance，default funcName>
// @import *<Path, required> <Alias>
// @injectParam *<ParamName, required> <Instance，default paramName>
// @injectRecv *<ParamName, required> <Instance，default paramName>
// @injectCtx *<ParamName, required>
```

> `@proxy` is used to indicate that the function supports dependency injection, and the system will automatically generate a proxy function with the same name as the function name, which can be accessed through a container object;
>
> `@InjectParam` is used for dependency injection of parameters;
>
> `@injectRecv` is used for dependency injection of the corresponding structure;
>
> `@injectCtx` is used to inject the container object itself;
>
> The parameters that have not been dependency injected will be retained in the generated proxy function;

### 2.4. WebApp Annotate (web server provided)
```
// @webAppProvide <WebApp，default WebApp>
// @static *<Path, required> *<Dirname, required> [Features: Compress|Download|Browse] <Index> <MaxAge>
```

> `@webAppProvide` Configure webApp. If not configured, the system will generate an instance named WebApp by default.
>
> If the webApp is not used in the code, the system will not generate a WebApp instance, which is also to adapt to non web projects.
>
> `@static` is used to configure static resource files, such as PNG, CSS, JS, HTML, etc


### 2.5. Router Annotate (like swag)
```
// @route *<Path, required> [Method: get|post]
// @webApp <WebApp，default WebApp>
// @injectCtx *<ParamName, required>
// @produce <Response Format: json | x-www-form-urlencoded | xml | plain | html | mpfd | json-api | json-stream | octet-stream | png | jpeg | gif>
// @param *<ParamName, required> <type:query|path|header|body|formData> <DataType> <IsRequired> <Description>
```

> `@router` will cause the system to generate a proxy function with the same name as the function to complete parameter parsing and injection, which can also be changed through `@proxy`;
>
> `@webApp` is used to associate webApp instances. WebApp instance is provided by `@webAppProvide` and the default instance name is "WebApp";
>
> `@injectWebCtx` is used to inject the webCtx of the current request, and can only be used for `@router` and `@middleware`;
>
> `@product` is used to define the return data type and can only be used for `@router`;
>
> `@param` is used for parsing request parameters and supports the following formats:
> * query: Get param, such as "/index.html?a=1";
> * path: router param, such as "/article/:id";
> * header: header param;
> * body: Body binary stream, note that only one body parameter can be defined;
> * formData: multipart/form format param;
>
> <b>Note: `@router` requires that each parameter must be configured with dependency injection</b>


### 2.6. Middleware Annotate
```
// @middleware *<Path, required>
// @webApp <WebApp，default WebApp>
// @injectCtx *<ParamName, required>
// @param *<ParamName, required> <type:query|path|header|body|formData> <DataType> <IsRequired> <Description>
```

> `@middleware` will cause the system to generate a proxy function with the same name as the function to complete parameter parsing and injection, which can also be changed through `@proxy`;
>
> <b>Note: `@middleware` requires that each parameter must be configured with dependency injection</b>

## 3. Generated code

### 3.1. PreConstruct & PostConstruct
Original structure
```go
// WebApp
// @provide WebCtxAlias
// @injectField Database
// @preConstruct model.PrepareWebCtxAlias
// @postConstruct model.WebCtxAliasLoaded
type WebApp struct {
    // @inject
    *Config
    *Database
    MiddleWares []*MiddleWare
    Routers     []*Router
}
```
> preConstruct function, Requirement must be no parameter and return type must be the same as the original struct.
```go
func PrepareWebCtxAlias() *WebApp {
// preConstruct function
    return &WebApp{}
}
```
> postConstruct function, Requirement the parameter must be of the original structure type.

```go
func WebCtxAliasLoaded(webApp *WebApp) {
// postConstruct function
}
```

### 3.2. Proxy Function

Original Function
```go
// WebCtxAliasLoaded
// @proxy
// @injectParam database Database
// @injectParam ctx
func WebCtxAliasLoaded(ctx *ctx.Ctx/*Special inject*/, webApp *WebApp/*uninjected*/, database *Database/*param inject*/) {
	fmt.Printf("WebCtxAliasLoaded: %v\n%v\n", webApp, database)
	ctx.TestLogin(webApp)
}
```
Generated Function
```go
func (ctx *Ctx/*Same name function in Ctx*/) WebCtxAliasLoaded(WebApp *model.WebApp/*Keep uninjected*/) {
	model.WebCtxAliasLoaded(ctx/*Special inject*/, WebApp, ctx.Database/*param inject*/)
}
```

> Many times, certain functions have strict formats, such as `@preConstruct`,`@postConstruct`, etc. 
> However, we also hope to rely on injecting some additional parameters to facilitate subsequent business operations. 
> At this time, we can consider using `@proxy`.
> 
> The function generated by `@proxy` only retains uninjected parameters, which can be used to generate a function that conforms to a certain format.
>
> Taking the postConstruct function as an example, it requires a structural type parameter, so we only need to keep one structural type parameter in the original function and not inject it.
> 
> Internally, the proxy function will assemble the required data based on injection rules and ultimately call the original function.
> 
> At this point, the original function obtains the injected parameters.
>

### 3.3. Web generated
Annotation configuration
```go
// ConfigureWebApp
// @webAppProvide WebApp
// @proxy WebAppStartup1
// @injectParam config Config
// @static /images /images
// @static /css /css [Compress,Browse]
// @static /js /js [Compress,Download,Browse] index.html 86400
func ConfigureWebApp(config *model.Config, extraParam int) (string, uint, error) {
	return config.Host, config.Port, nil
}

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

// LoginController
// @router /api/login [post]
// @param username query string true username
// @param password query string true password
func LoginController(username string, password string) error {
return nil
}
```
Generating
```go
func (ctx *Ctx) WebAppStartup1(extraParam int) error {
    ctx.WebApp.Static("/images", "")
    ctx.WebApp.Static("/css", "", fiber.Static{Compress: true, Browse: true})
    ctx.WebApp.Static("/js", "", fiber.Static{Compress: true, Download: true, Browse: true, Index: "index.html", MaxAge: 86400})
    /// middleware register
	ctx.WebApp.Group("/api", ctx.CorsMiddleware)
    /// router register
    ctx.WebApp.Post("/api/login", ctx.LoginController)
    host, port, err := init.ConfigureWebApp(webApp, ctx.Config, extraParam)
    if err != nil {
        return err
    }
    return ctx.WebApp.Listen(fmt.Sprintf("%s:%d", host, port))
}

func (ctx *Ctx) CorsMiddleware(webCtx *fiber.Ctx) (err error) {
    body := &model.Config{}
    err = BodyParser(webCtx, body)
    if err != nil {
        return err
    }
    header := Header(webCtx, "header")
    queryBool, err := QueryBool(webCtx, "queryBool")
    if err != nil {
        return err
    }
    paramsInt, err := ParamsInt(webCtx, "paramsInt")
    if err != nil {
        return err
    }
    formFloat, err := FormFloat(webCtx, "formFloat")
    if err != nil {
        return err
    }
    return middleware.CorsMiddleware(ctx, webCtx, body, header, paramsInt, queryBool, formFloat)
}

func (ctx *Ctx) LoginController(webCtx *fiber.Ctx) (err error) {
    username := Query(webCtx, "username")
    password := Query(webCtx, "password")
    return router.LoginController(username, password)
}
```

Start web service
```go
func main() {
    err := ctx.New().WebAppStartup1(100)
    if err != nil {
        return
    }
}
```

> `@webAppProvide` is used to define web application instances. If not configured, the system will create an instance named 'WebApp' and a startup function named 'WebAppStartup' by default;
>
> The instance name is modified through `@webAppProvide`, and the startup function name is modified through `@proxy`;
>
> `@router` and `@middleware` are associated with an instance through `@webApp`, which defaults to the instance name associated with 'WebApp';
>
> The instance names of `@webAppProvide` and `@webApp` must be consistent to ensure their association;

> <b>Note: `@middleware` and `@router` require that each parameter must be configured with dependency injection The system will automatically create code functions that conform to the webApp call format.</b>
> 
> Only `@middleware` and `@router` can inject `@webCtx`, which represents the context of this request;
> 
> For more tips on using webCtx, you can read ` * fiber Ctx ` related documents.
> 

### 3.4. Generated Directory
```
/ctx
    |- gen_ctx.go
        --------------------------------
        # gen segment: Struct #
        --------------------------------
        type Ctx struct {
            {{range SingletonInstances}}
            {{Instance}} {{Name}}
            {{end}}
            {{range WebAppInstances}}
            {{WebApp}} *fiber.App
            {{end}}
        }
      
        -----------------------------------
        # gen segment: Singleton instance #
        -----------------------------------
        func New() *Ctx {
            ctx := &Ctx{}
            {{range SingletonInstances}}
                {{if PreConstruct}}
                    ctx.{{Instance}} := {{PreConstruct}}()
                {{else}}
                    ctx.{{Instance}} = &{{Package}}.{{Name}}{}
                {{end}}
            {{end}}
            {{range WebAppInstances}}
                ctx.{{WebApp}} := fiber.New()
            {{end}}
            
            {{range SingletonInstances}}
                {{range InjectFields}}
                    {{if FieldInstance == "Ctx"}}
                    ctx.{{Instance}}.{{FieldInstance}} = ctx
                    {{else}}
                        {{if IsSingleton}}
                        ctx.{{Instance}}.{{FieldInstance}} = ctx.{{StructInstance}}
                        {{else if IsMultiple}}
                        ctx.{{Instance}}.{{FieldInstance}} = ctx.New{{StructInstance}}()
                        {{end}}
                    {{end}}
                {{end}}
            {{end}}
            
            {{range SingletonInstances}}
                {{if PostConstruct}}
                    {{PostConstruct}}(
                        ctx.{{Instance}},
                    )
                {{end}}
            {{end}}
            return ctx
        }
    |- gen_constructor.go
        ------------------------------------
        # gen segment: Multiple instance #
        ------------------------------------
        {{range MultipleInstances}}
            func (ctx *Ctx) New{{Instance}}() *{{Type}} {
                {{if PreConstruct}}
                    {{Instance}} := {{PreConstruct}}()
                {{else}}
                    {{Instance}} := &{{Package}}.{{Name}}{}
                {{end}}
                {{range Fields}}
                    {{if IsInject}}
                        {{if FieldInstance == "Ctx"}}
                        {{Instance}}.{{FieldName}} = ctx
                        {{else}}
                            {{if IsSingleton}}
                            {{Instance}}.{{FieldName}} = ctx.{{FieldInstance}}
                            {{else if IsMultiple}}
                            {{Instance}}.{{FieldName}} = ctx.New{{FieldInstance}}()
                            {{end}}
                        {{end}}
                    {{end}}
                {{end}}
                
                {{if PostConstruct}}
                    {{PostConstruct}}(
                        {{Instance}},
                    )
                {{end}}
                return {{Instance}}
            }
        {{end}}
        
    |- gen_func.go
        ------------------------------------
        # gen segment: Func inject #
        ------------------------------------
        {{range FuncInstances}}
            func (ctx *Ctx) {{Proxy}}(
                {{range NormalParams}}
                {{ParamInstance}} {{ParamType}},
                {{end}}
            ) (
                {{range Results}}
                {{ResultName}} {{ResultType}},
                {{end}}
            ) {
                return {{Package}}.{{FuncName}}(
                    {{range Params}}
                        {{if IsInject}}
                            {{if ParamInstance == "Ctx"}}
                            ctx,
                            {{else}}
                                {{if IsSingleton}}
                                ctx.{{ParamInstance}},
                                {{else if IsMultiple}}
                                ctx.New{{ParamInstance}}(),
                                {{end}}
                            {{end}}
                        {{else}}
                            {{ParamInstance}},
                        {{end}}
                    {{end}}
                )
            }
        {{end}}
        
    |- gen_method.go
        ------------------------------------
        # gen segment: Method inject #
        ------------------------------------
        func (ctx *Ctx) {{Proxy}}(
            {{if !Recv.IsInject}}
            {{Recv.Name}} {{Recv.Type}},
            {{end}}
            {{range NormalParams}}
            {{ParamInstance}} {{ParamType}},
            {{end}}
        ) (
            {{range Results}}
            {{ResultName}} {{ResultType}},
            {{end}}
        ) {
            return {{Recv.Name}}.{{FuncName}}(
                {{if !Recv.IsInject}}
                {{Recv.Name}},
                {{end}}
                {{range Params}}
                    {{if IsInject}}
                        {{if ParamInstance == "Ctx"}}
                        ctx,
                        {{else}}
                            {{if IsSingleton}}
                            ctx.{{ParamInstance}},
                            {{else if IsMultiple}}
                            ctx.New{{ParamInstance}}(),
                            {{end}}
                        {{end}}
                    {{else}}
                        {{ParamInstance}},
                    {{end}}
                {{end}}
            )
        }
    |- gen_web.go
        ------------------------------------
        # gen segment: WebApp startup #
        ------------------------------------
        {{range WebAppInstances}}
        func (ctx *Ctx) {{Proxy |default "WebAppStartup"}}(
            {{if FuncName != ""}}
                {{if !Recv.IsInject}}
                {{Recv.Name}} {{Recv.Type}},
                {{end}}
                {{range NormalParams}}
                {{ParamInstance}} {{ParamType}},
                {{end}}
            {{else}}
                host string,
                port uint,
            {{end}}
        ) error {
            {{range Statics}}
            ctx.{{WebApp}}.Static({{Path}}, {{Dirname}}, 
                {{if Features || Index || MaxAge}}
                fiber.Static{
                    {{range Features}}
                    {{Feature}}: true,
                    {{end}}
                    {{if Index}}
                    Index: {{Index}},
                    {{end}}
                    {{if MaxAge}}
                    MaxAge: {{MaxAge}},
                    {{end}}
                }
                {{end}}
            )
            {{end}}
            {{range Middlewares}}
            ctx.{{WebApp}}.Group({{Path}}, ctx.{{Proxy}})
            {{end}}
            
            {{range Routers}}
                {{range Methods}}
                ctx.{{WebApp}}.{{Method}}({{Path}}, ctx.{{Proxy}})
                {{end}}
            {{end}}
            
            {{if FuncName != ""}}
            host, port, err := {{Package}}.{{FuncName}}(
                {{range Params}}
                    {{if param.Source == "ctx"}}
                        ctx,
                    {{else if param.Source == "inject"}}
                        ctx.{{ParamInstance}},
                    {{else}}
                        {{ParamInstance}},
                    {{end}}
                {{end}}
            )
            if err != nil {
                return err
            }
            {{end}}
            return webApp.Listen(fmt.Sprintf("%s:%d", host, port))
        }
        {{end}}
        
        ------------------------------------
        # gen segment: Middleware #
        ------------------------------------
        {{range WebAppInstances}}
            {{range Middlewares}}
            func (ctx *Ctx) {{Proxy}}(
                webCtx *fiber.Ctx,
                {{if !Recv.IsInject}}
                {{Recv.Name}} {{Recv.Type}},
                {{end}}
                {{range NormalParams}}
                {{ParamInstance}} {{ParamType}},
                {{end}}
            ) (err error) {
                {{if BodyParam}}
                    {{if Type == "[]btye"}}
                        {{ParamInstance}} := Body(webCtx)
                    {{else if Type == "string"}}
                        {{ParamInstance}} := BodyString(webCtx)
                    {{else}}
                        {{ParamInstance}} := &{{Package}}.{{ParamType}}{}
                        err := BodyParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := HeaderParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := Header(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := HeaderInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := HeaderBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := HeaderFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := HeaderParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := QueryParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := Query(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := QueryInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := QueryBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := QueryFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := QueryParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := FormParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := FormString(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := FormInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := FormBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := FormFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "file"}}
                        {{ParamInstance}}, err := FormFile(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := FormParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                return {{Package}}.{{FuncName}}(
                    {{range _, param := Params}}
                        {{if param.Source == "ctx"}}
                            ctx,
                        {{else if param.Source == "webCtx"}}
                            webCtx,
                        {{else if param.Source == "inject"}}
                            ctx.{{ParamInstance}},
                        {{else}}
                            {{param.Instance}},
                        {{end}}
                    {{end}}
                )
            }
            {{end}}
        {{end}}
        
        ------------------------------------
        # gen segment: Router #
        ------------------------------------
        {{range WebAppInstances}}
            {{range Routers}}
            func (ctx *Ctx) {{Proxy}}(
                webCtx *fiber.Ctx,
                {{if !Recv.IsInject}}
                {{Recv.Name}} {{Recv.Type}},
                {{end}}
                {{range NormalParams}}
                {{ParamInstance}} {{ParamType}},
                {{end}}
            ) (err error) {
                {{if BodyParam}}
                    {{if Type == "[]btye"}}
                        {{ParamInstance}} := Body(webCtx)
                    {{else if Type == "string"}}
                        {{ParamInstance}} := BodyString(webCtx)
                    {{else}}
                        {{ParamInstance}} := &{{Package}}.{{ParamType}}{}
                        err := BodyParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := HeaderParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := Header(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := HeaderInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := HeaderBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := HeaderFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := HeaderParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := QueryParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := Query(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := QueryInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := QueryBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := QueryFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := QueryParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := FormParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := FormString(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := FormInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := FormBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := FormFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "file"}}
                        {{ParamInstance}}, err := FormFile(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := FormParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                return {{Package}}.{{FuncName}}(
                    {{range _, param := Params}}
                        {{if param.Source == "ctx"}}
                            ctx,
                        {{else if param.Source == "webCtx"}}
                            webCtx,
                        {{else if param.Source == "inject"}}
                            ctx.{{ParamInstance}},
                        {{else}}
                            {{param.Instance}},
                        {{end}}
                    {{end}}
                )
            }
            {{end}}
        {{end}}
```
