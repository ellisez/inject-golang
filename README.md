# inject-golang
Provide container of DI(Dependency Injection) for golang.

Language: English [中文](README_cn.md) 

### 1. Install and Operation

```shell
go install github.com/ellisez/inject-golang
```
### 1.1. Configuration generator
> By default, the current module is scanned and all annotations are generated.
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
> Use `-m` to specify that only partial annotations are generated, support setting multiple, separated by commas, and optional values such as 'singleton','multiple','fun', and 'web'

```go
//go:generate inject-golang -m singleton,multiple
func main() {
    ctx.New()
}
```

> Can specify scanning modules: supports multiple modules, defaults to the current module, external modules must be imported.
>
> Symbol "." represents current package, The system already supports `go.work`.

```go
//go:generate inject-golang -m singleton,web github.com/ellisez/inject-golang/examples-work .
func main() {
    ctx.New()
}
```


To learn more about commands, please run `inject-golang -h`.

### 1.2. Run generator
```shell
go generate -run inject-golang
```

### 1.3. Clean Code
```shell
inject-glang --clean
```

## 2. Annotate

### 2.1. Disable Annotations
Use `//@` Adding an exclamation mark at the beginning, which means adding an exclamation mark before the annotation, will not be recognized or parsed by the system.
```go
// !@proxy
```

### 2.2. Proxy Function Annotations

`Proxy Function Annotations` is to enable dependency injection on normal functions.
```
// @proxy <Instance，default funcName>
// @override
// @import *<Path, required> <Alias>
// @injectParam *<ParamName, required> <Instance，default paramName> <operator, ""|&|*|cast>
// @injectRecv *<ParamName, required> <Instance，default paramName>
// @injectCtx *<ParamName, required>
// @injectFunc *<ParamName, required> <Instance，default paramName> <operator, ""|call>
```

> `@proxy ` causes the system to generate a proxy function with the same name as the original function by default.
> Proxy functions can be accessed through container objects
>
> `@override` indicates support for overloading, and when encountering instances with the same name, the latter will overwrite the former; The default is to disable overloading, and an error will be reported when the same name is used.
> 
> `@injectParam` is used for dependency injection of parameters;
>
> `@injectRecv` is used for structural dependency injection of member functions;
>
> `@injectCtx` is used to inject the container object itself;
>
> `@injectFunc` is used to inject parameters of function types;
> 
> `@injectParam` supports type conversion, where `&` represents the address of the value, `*` represents the value corresponding to the address, `cast` represents strong type conversion, and the default value is `""` indicating no conversion.

>
> Note: The parameters that have not been dependency injected will be retained in the generated proxy function;


### 2.3. Instance Annotations
`Instance annotations` refers to declaring an instance by marking annotations with a constructor
```
// @provide <Instance, default ReturnType> <singleton default|multiple|argument> <type, default ReturnType>
// @override
// @order <Instance creation order, numbers or strings>
// @import *<Path, required> <Alias>
// @injectParam *<ParamName, required> <Instance，default paramName> <operator, ""|&|*|cast>
// @injectCtx *<ParamName, required>
// @injectFunc *<ParamName, required> <Instance，default paramName> <operator, ""|call>
// @handler *<called after creation, required>
```

> The system will load the container package and the package where the annotation is located by default, but if there are additional packages, they need to be declared with `@import`;
>
> The constructor marked by `@provide` must have and only have one return type, support dependency injection, but must inject each parameter.
> 
> `@provide` requires `@order` to define the creation order to prevent instances that have not yet been initialized from being injected;
>
> `@provide` modes include: `singleton` represents global uniqueness, `multiple` represents the ability to create multiple, and `argument` represents only existing during the startup process;
> 
> `@override` indicates support for overloading, and when encountering instances with the same name, the latter will overwrite the former; The default is to disable overloading, and an error will be reported when the same name is used.
> 
> The function pointed to by `@handler` must have no parameters, and it will be called after the instance is created.
> 
> `@handler` can carry the package name, for example: `*model.Database`, which represents calling the original function.But if it does not include a package name, it is a proxy function.
>
> We generally do not recommend using the original function directly, but rather its proxy function, which can help us expand the parameters of other dependency injections;
> Proxy function usage, please refer to [Func Annotate (use for all func)](#23-func-annotate-use-for-all-func)
>
>


### 2.4. WebApp Annotates (web server provided)
```
// @webProvide <instance，default WebApp>
// @static *<Path, required> *<Dirname, required> [Features: Compress|Download|Browse] <Index> <MaxAge>
```

> `@webProvide` Configure Web Application. If not configured, the system will generate an instance named `WebApp` by default.
>
> If the Web Application is not used in the code, the system will not generate a Web Application instance, which is also to adapt to non web projects.
>
> The startup function for web applications, in the format of `instance + Startup`, defaults to `WebAppStartup`.
> 
> The original function marked by `@webProvide` must return three parameters: `host`, `port`, and `err`.
> 
> `@static` is used to configure static resource files, such as PNG, CSS, JS, HTML, etc


### 2.5. Router Annotates (like swag)
```
// @route *<Path, required> [Method: get|post]
// @webApp <WebApp，default WebApp>
// @injectCtx *<ParamName, required>
// @produce <Response Format: json | x-www-form-urlencoded | xml | plain | html | mpfd | json-api | json-stream | octet-stream | png | jpeg | gif>
// @param *<ParamName, required> <type:query|path|header|body|formData> <DataType> <IsRequired> <Description>
```

> `@router` will cause the system to generate a proxy function with the same name as the function to complete parameter parsing and injection, which can also be changed through `@proxy`;
>
> `@webApp` is used to associate webApp instances. WebApp instance is provided by `@webProvide` and the default instance name is "WebApp";
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


### 2.6. Middleware Annotations
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

### 3.1. Create an instance
Annotations on Construct Functions
```go
// PrepareServerAlias example for proxy handler
// @provide ServerAlias _ model.ServerInterface
// @order "step 4: Setting Server"
// @import github.com/ellisez/inject-golang/examples/model
// @injectParam config
// @injectParam database
// @handler ServerAliasLoaded
func PrepareServerAlias(config *model.Config, database *model.Database) *model.Server {
    fmt.Println("call WebAppAlias.PrepareWebAppAlias")
    return &model.Server{
    Config:   config,
        Database: database,
    }
}
```
> Constructor, must have only one return type.
> 
> In the example, the instance type specified `model.ServerInterface` is an interface type, while the actual creation type is the return type `*model.Server`.

`@handle` Annotations on the original function
```go
// ServerAliasLoaded example for injection proxy
// @proxy
// @import "github.com/ellisez/inject-golang/examples/model"
// @injectParam database Database
// @injectCtx appCtx
// @injectParam server ServerAlias cast
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
```

> The handler function requires no parameters. In the above example, a proxy function was used because the original function needs to inject other instances to perform specific functions.

The proxy function generated by `@handler` is as follows:
```go
// Generate by annotations from handler.ServerAliasLoaded
func (ctx *Ctx) ServerAliasLoaded() {
    handler.ServerAliasLoaded(ctx, ctx.serverAlias.(*model.Server), ctx.database, &ctx.isReady, ctx.NewEvent(), ctx.NewListener())
}
```
> The generated proxy function will retain the parameters that have not been injected, and only when all parameters are injected can a parameterized proxy function be generated, which is the format requirement for `handler`.
> 
> The parameter `server` uses a `cast` operator, so strong conversion was automatically performed in the generated code `cast` is commonly used for converting interfaces to structures.
> 
> The parameter `isReady` uses a `&` operator, so the address retrieval operation is automatically performed in the generated code. 
> `IsReady` is a basic type that takes a pointer to allow it to be modified.

### 3.2. Web Generated Codes

Annotation configuration
```go
// ConfigureWebApp
// @webProvide instance
// @import github.com/ellisez/inject-golang/examples/model
// @proxy WebAppStartup1
// @injectParam config Config
// @static /images /images
// @static /css /css [Compress,Browse]
// @static /js /js [Compress,Download,Browse] index.html 86400
func ConfigureWebApp(config *model.Config, defaultPort uint) (string, uint, error) {
    if config.Port == 0 {
        defaultPort = config.Port
    }
    return config.Host, defaultPort, nil
}

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
    body *model.Config,
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

// LoginController
// @router /api/login [post]
// @import github.com/ellisez/inject-golang/examples/model
// @param username query string true 
// @param password query string true 
// @injectParam server ServerAlias
func LoginController(username string, password string, server *model.Server) error {
    fmt.Printf("call LoginController: %s, %s\n", username, password)
    server.TriggerEvent("login", map[string]any{
        "username": username,
        "password": password,
    })
    return nil
}
```
Generating
```go
// Generate by annotations from provide.ConfigureWebApp
func (ctx *Ctx) WebAppStartup(defaultPort uint) error {
    ctx.WebApp().Static("/images", "/images")
    ctx.WebApp().Static("/css", "/css", fiber.Static{Compress: true, Browse: true})
    ctx.WebApp().Static("/js", "/js", fiber.Static{Compress: true, Download: true, Browse: true, Index: "index.html", MaxAge: 86400})
    ctx.WebApp().Group("/api", ctx.CorsMiddleware)
    ctx.WebApp().Post("/api/register", ctx.RegisterController)
    ctx.WebApp().Post("/api/login", ctx.LoginController)
    host, port, err := provide.ConfigureWebApp(ctx.config, defaultPort)
    if err != nil {
        return err
    }
    return ctx.WebApp().Listen(fmt.Sprintf("%s:%d", host, port))
}

// Generate by annotations from middleware.CorsMiddleware
func (ctx *Ctx) CorsMiddleware(webCtx *fiber.Ctx) (err error) {
    body := &model.Config{}
    err = utils.BodyParser(webCtx, body)
    if err != nil {
        return err
    }
    header := utils.Header(webCtx, "header")
    queryBool, err := utils.QueryBool(webCtx, "queryBool")
    if err != nil {
        return err
    }
    paramsInt, err := utils.ParamsInt(webCtx, "paramsInt")
    if err != nil {
        return err
    }
    formFloat, err := utils.FormFloat(webCtx, "formFloat")
    if err != nil {
        return err
    }
    return middleware.CorsMiddleware(ctx, webCtx, body, header, paramsInt, queryBool, formFloat)
}

// Generate by annotations from router.LoginController
func (ctx *Ctx) LoginController(webCtx *fiber.Ctx) (err error) {
    username := utils.Query(webCtx, "username")
    password := utils.Query(webCtx, "password")
    return router.LoginController(username, password, ctx.serverAlias)
}
```
Calling the web's startup function
```go
func main() {
    c := factory.New()
    err := c.WebAppStartup(3001)
    if err != nil {
        return
    }
}
```

> `@webProvide` is used to define web application instances. If not configured, the system will create an instance named 'WebApp' and a startup function named 'WebAppStartup' by default;
> 
> The instance name is modified through `@webProvider`, and the startup function will also change accordingly;
>
> `@router` and `@middleware` are associated with an instance through `@webApp`, which defaults to the instance name associated with 'WebApp';
>
> The instance names of `@webProvide` and `@webApp` must be consistent to ensure their association;


> <b>Note: `@middleware` and `@router` require that each parameter must be configured with dependency injection The system will automatically create code functions that conform to the webApp call format.</b>
>
> Only `@middleware` and `@router` can inject `@webCtx`, which represents the context of this request;
>
> For more tips on using webCtx, you can read ` * fiber Ctx ` related documents.
>

### 3.4. Directory structure
```
/----
    /ctx                    Generate code, please do not modify it
        |- gen_ctx.go       ctx interface
    /model                  Defining Structures
    /provide                Instance annotations
    /middleware             Web middleware annotations
    /router                 Web router annotations
    /controller             Processing functions for web requests
    /service                Processing functions for services and transactions
    /handler                Unclassified processing function
    |- main.go              Startup function
```

> The `global` directory is not recommended, and ctx should be used to access global variables to avoid code bloating.
>
> It is not recommended to write a large amount of startup sequence code for `main.go`. Instead, it should be completed elegantly through the `@order` sequence of the instance

## 4.Usage specifications

### 4.1 Package Name Specification

> To improve the readability of the code, please remember to keep the package name the same as the directory name except for special types of packages, such as version packages and application packages.

The following is the naming convention for special type packages:
* `Version package`: The directory name format is `v[\d.]+`, with the letter v followed by a number; 
<br/>It represents a package within a specified version range;
<br/>The package name of the version package should be the previous directory name, such as `fiber` in `github.com/gofiber/fiber/v2`;

* `Application package`: The package is named main and serves as the program startup entry point, usually appearing in the root directory of the mod module.
<br/>Golang stipulates that main cannot be imported, so even if global variables are defined in the main package, it cannot be accessed by other packages;
<br/>Although the application package cannot be imported, the system still read the annotations inside the package;

### 4.2 Circular dependency problem

Golang prohibits two packages from importing each other. 

> To avoid this, we should adhere to the `declaration and call separation` principle in design.

The specific actions are as follows:
* Two types of packages should be prepared, one for declaration and the other for calling; 
<br/>Call packages can import dependent import declaration packages, but declaration packages prohibit the import of call packages;

* The `declaration package` should include the structures of annotation declarations such as `@provide`, `@webProvide`, which can provide the rules for creating instances; <br/>Recommend the package name as `model`;
* The `call package` should include the functions of annotation declarations such as `@handler`, `@proxy`, `@middleware`, and `@router`, which can provide function callbacks for dependency injection; <br/>Recommended package name is `handler`;



> Note: The declaration package prohibits the use of the `@injectCtx` annotation.
>
> For example, when declaring the member function of a package using '@ proxy', the context of 'Ctx' should not be injected, but should be replaced with value injection and function injection.
>
> Value injection includes `@injectParam` and `@injectRecv`, while function injection includes `@injectFunc`.

### 4.3. Clear Annotate Override
In the cross module engineering of `go.work`, many times the instances and proxy functions generated by annotations defined by submodules are not compatible with the main module. 

The traditional approach can only be to create new instances and proxy functions, or to expand the code of submodules.

If new code is written, the old code of the submodule will become invalid and cause redundancy, while expanding the submodule will cause the submodule and the main module to call each other. Obviously, 
both of these methods are not ideal.

This is our recommendation to use the `@override` annotation to overwrite.

Annotation coverage can achieve the same instance name, while the latter covers the former, and the former's annotations do not leave any traces in the generated code.

For example, for a structural instance defined by a submodule, the main module needs to expand some fields on the structure. 

In this case, only annotations are needed to overwrite the type of the instance.
