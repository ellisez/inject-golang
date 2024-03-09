# inject-golang
提供golang版依赖注入容器

语言: [English](README.md) 中文

### 1. 安装与运行

```shell
go install github.com/ellisez/inject-golang
```
### 1.1. 配置生成器
> 默认扫描当前模块并生成所有注解
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
> 使用`-m`可指定只生成部分注解, 支持设置多个, 用逗号分隔, 可选singleton, multiple, func, web

```go
//go:generate inject-golang -m singleton,multiple
func main() {
    ctx.New()
}
```

> 可指定扫描模块: 支持多个, 默认为当前模块, 外部模块则必须要被引入.
> 
> "." 表示当前包, 系统已支持go.work引入方式.

```go
//go:generate inject-golang -m singleton,web github.com/ellisez/inject-golang/examples-work .
func main() {
    ctx.New()
}
```

想了解更多配置命令, 请运行`inject-golang -h`.

### 1.2. 运行生成器
```shell
go generate -run inject-golang
```

### 1.3. 清空代码
```shell
inject-glang --clean
```

## 2. 注解
### 2.1. 结构体上的注解：
```
// @provide <实例名，默认同类名> <singleton默认|multiple>
// @import *<模块加载路径，必填> <模块名>
// @injectField *<字段名，必填> <实例名，默认同类名> <读取函数Getter, 仅在私有属性时中作> <写入函数Setter, 仅在私有属性时工作>
// @preConstruct *<构造前调用函数，必填>
// @postConstruct *<构造后调用函数，必填>
```

> 系统默认会加载容器包以及注解所在代码的包, 但如果有额外的包, 则需要使用`@import`声明.
> 
> `@preConstruct`是一个无参但返回结构体类型的函数, 主要用于创建前给定结构体的初始值, 不支持注入依赖;
> 
> `@postConstruct`是一个唯一参数为结构体类型的函数, 主要用于获取已经创建完成的结构体实例;
> 
> `@preConstruct`与`@postConstruct`可携带包名, 但需要`@import`来引入, 如"*model.Database".
> 
> 一般我们不推荐直接使用原始函数, 而是使用它的代理函数, 这样能帮我们扩充其他依赖注入的参数; 
> 代理函数用法, 请参照[方法上注解](#23-方法上的注解-适用于所有方法)
>
>

### 2.2. 结构体的属性注解：
```
// @inject <实例名，默认同类名> <读取函数Getter, 仅在私有属性时中作> <写入函数Setter, 仅在私有属性时工作>
```
### 2.3. 方法上的注解 (适用于所有方法)
```
// @proxy <代理方法名，默认同方法名>
// @import *<模块加载路径，必填> <模块名>
// @injectParam *<参数名，必填> <实例名，默认同类名>
// @injectRecv *<参数名，必填> <实例名，默认同类名>
// @injectCtx *<参数名, 必填>
```

> `@proxy`用于标记该函数支持依赖注入, 最终系统会自动生成一个与函数名同名代理函数, 可通过容器对象访问代理函数.
> 
> `@injectParam`用于参数的依赖注入;
> 
> `@injectRecv`用于所属结构体的依赖注入;
> 
> `@injectCtx`用于注入容器对象本身;
> 
> 未被依赖注入的参数则会保留到生成代理函数中;

### 2.4. WebApp注解 (提供了web服务器)
```
// @webAppProvide <WebApp，默认名为WebApp>
// @static *<访问路径，必填> *<匹配目录，必填> [特征: Compress|Download|Browse] <目录的Index文件> <过期时间MaxAge>
```

> `@webAppProvide`配置webApp, 不进行配置时系统默认会生成一个名为WebApp的实例. 
> 
> 如果webApp未被代码中使用, 则系统不会生成WebApp实例, 这也是为了适配与非web项目.
> 
> `@static`用于配置静态资源文件, 如png,css,js,html等

### 2.5. 路由方法上的注解（参照swag注解）：
```
// @router *<Path必填> [Method: get|post]
// @webApp <WebApp，默认名为WebApp>
// @injectWebCtx *<参数名, 必填>
// @produce <返回格式: json | x-www-form-urlencoded | xml | plain | html | mpfd | json-api | json-stream | octet-stream | png | jpeg | gif>
// @param *<参数名，必填> *<取值类型，必填:query|path|header|body|formData> <接收类型> <必填与否> <参数说明>
```

> `@router`会让系统生成一个与函数同名的代理函数, 以完成参数的解析和注入, 也可以通过`@proxy`更改.
> 
> `@webApp`用于关联webApp实例, webApp由`@webAppProvide`提供, 默认实例名为"WebApp".
> 
> `@injectWebCtx`用于注入当前请求的webCtx, 只能用于`@router`和`@middleware`;
> 
> `@produce`用于定义返回数据类型, 只能用于`@router`;
> 
> `@param`用于请求参数的解析, 支持下列格式:
> * query: Get参数, 如/index.html?a=1;
> * path: 路由参数, 如/article/:id;
> * header: 头部信息参数;
> * body: body二进制流, 注意只能由一个body参数;
> * formData: multipart/form方式提交的数据;
> 
> <b>注意: `@router`要求每个参数必须都被配置依赖注入</b>

### 2.6. 路由中间件上的注解:
```
// @middleware *<Path必填>
// @webApp <WebApp，默认名为WebApp>
// @injectWebCtx *<参数名, 必填>
// @param *<参数名，必填> *<参数类型，必填:query|path|header|body|formData> <接收类型> <必填与否> <参数说明>
```

> `@middleware`会让系统生成一个与函数同名的代理函数, 以完成参数的解析和注入, 也可以通过`@proxy`更改.
>
> <b>注意: `@middleware`要求每个参数必须都被配置依赖注入</b>

## 3. 生成模板

### 3.1. PreConstruct & PostConstruct
原结构体
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
> preConstruct函数, 要求必须是无参并且返回类型必须和结构体类型一样

```go
func PrepareWebCtxAlias() *WebApp {
// preConstruct function
    return &WebApp{}
}
```

> postConstruct函数, 要求参数必须是结构体类型

```go
func WebCtxAliasLoaded(webApp *WebApp) {
// postConstruct function
}
```

### 3.2. 方法代理

原函数
```go
// WebCtxAliasLoaded
// @proxy
// @injectParam database Database
// @injectParam ctx
func WebCtxAliasLoaded(ctx *ctx.Ctx/*特殊注入*/, webApp *WebApp/*未注入*/, database *Database/*属性注入*/) {
	fmt.Printf("WebCtxAliasLoaded: %v\n%v\n", webApp, database)
	ctx.TestLogin(webApp)
}
```
生成的目标函数
```go
func (ctx *Ctx/*属于容器的同名函数*/) WebCtxAliasLoaded(WebApp *model.WebApp/*保留未注入*/) {
	model.WebCtxAliasLoaded(ctx/*特殊注入*/, WebApp, ctx.Database/*属性注入*/)
}
```

> 许多时候, 某些函数有严格的格式, 如`@preConstruct`, `@postConstruct`等, 但我们又希望能依赖注入一些额外的参数, 以方便完成后续业务, 这个时候我们就可以考虑使用`@proxy`.
> 
> `@proxy`生成的代理函数, 只会保留未注入的参数, 利用这一点就能够生成符合一定格式的函数.
> 
> 以postConstruct函数为例, 它需要一个结构体类型的参数, 所以我们只需要在原函数中保留一个结构体类型参数让它不进行注入即可.
> 
> 而代理函数内部又会根据注入规则, 组装所需数据并最终调用原函数.
> 
> 此时原函数就获得了, 已经被注入后的参数.


### 3.3. web生成代码
注解配置
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
// @param username query string true 用户名
// @param password query string true 密码
func LoginController(username string, password string) error {
return nil
}
```
生成代码
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
外部启动web
```go
func main() {
    err := ctx.New().WebAppStartup1(100)
    if err != nil {
        return
    }
}
```

> `@webAppProvide`用于定义web应用实例，不配置时系统默认会创建名为`WebApp`的实例，以及名为`WebAppStartup`的启动函数,
>
> 实例名通过`@webAppProvide`修改, 启动函数名通过`@proxy`来修改.
>
> `@router`与`@middleware`通过`@webApp`关联实例, 默认关联`WebApp`的实例名.
>
> `@webAppProvide`与`@webApp`实例名必须保持一致, 才能确保关联关系.

> <b>注意: `@middleware`与`@router`要求每个参数都必须都要配置依赖注入. 系统会自动创建出符合webApp调用格式的代码函数</b>
>
> 只有`@middleware`和`@router`可以注入`@webCtx`, 它表示本次请求上下文.
>
> webCtx更多使用技巧, 可阅读`*fiber.Ctx`相关文档
>

### 3.4. 目录结构
```
/ctx
    |- gen_singleton.go
        --------------------------------
        # gen segment: Struct #
        --------------------------------
        type Ctx struct {
            {{range SingletonInstances}}
            {{PrivateName}} {{Name}}
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
                    ctx.{{PrivateName}} := {{PreConstruct}}()
                {{else}}
                    ctx.{{PrivateName}} = &{{Package}}.{{Name}}{}
                {{end}}
            {{end}}
            {{range WebAppInstances}}
                ctx.{{PrivateName}} := fiber.New()
            {{end}}
            
            {{range SingletonInstances}}
                {{range Fields}}
                    {{if Field.Source == "ctx"}}
                    ctx.{{PrivateName}}.{{FieldInstance}} = ctx
                    {{else if Field.Source == "inject"}}
                        {{if IsSingleton}}
                            {{if IsPrivate}}
                            ctx.{{PrivateName}}.{{Field.Setter}}(ctx.{{Field.Instance}}())
                            {{else}}
                            ctx.{{PrivateName}}.{{Field.Name}} = ctx.{{Field.Instance}}()
                            {{end}}
                        {{else if IsMultiple}}
                            {{if IsPrivate}}
                            ctx.{{PrivateName}}.{{Field.Setter}}(ctx.New{{Field.Instance}}())
                            {{else}}
                            ctx.{{PrivateName}}.{{Field.Name}} = ctx.New{{Field.Instance}}()
                            {{end}}
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
    |- gen_multiple.go
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
                            {{Instance}}.{{FieldName}} = ctx.{{FieldInstance}}()
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
                                ctx.{{ParamInstance}}(),
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
                            ctx.{{ParamInstance}}(),
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
            ctx.{{WebApp}}().Static({{Path}}, {{Dirname}}, 
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
            ctx.{{WebApp}}().Group({{Path}}, ctx.{{Proxy}})
            {{end}}
            
            {{range Routers}}
                {{range Methods}}
                ctx.{{WebApp}}().{{Method}}({{Path}}, ctx.{{Proxy}})
                {{end}}
            {{end}}
            
            {{if FuncName != ""}}
            host, port, err := {{Package}}.{{FuncName}}(
                {{range Params}}
                    {{if param.Source == "ctx"}}
                        ctx,
                    {{else if param.Source == "inject"}}
                        {{if IsSingleton}}
                        ctx.{{ParamInstance}}(),
                        {{else if IsMultiple}}
                        ctx.New{{ParamInstance}}(),
                        {{end}}
                    {{else}}
                        {{ParamInstance}},
                    {{end}}
                {{end}}
            )
            if err != nil {
                return err
            }
            {{end}}
            return ctx.{{WebApp}}().Listen(fmt.Sprintf("%s:%d", host, port))
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
                        {{ParamInstance}} := utils.Body(webCtx)
                    {{else if Type == "string"}}
                        {{ParamInstance}} := utils.BodyString(webCtx)
                    {{else}}
                        {{ParamInstance}} := &{{Package}}.{{ParamType}}{}
                        err := utils.BodyParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := HeaderParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := utils.Header(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := utils.HeaderInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := utils.HeaderBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := utils.HeaderFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := utils.HeaderParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := QueryParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := utils.Query(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := utils.QueryInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := utils.QueryBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := utils.QueryFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := utils.QueryParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := PathParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := utils.Params(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := utils.ParamsInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := utils.ParamsBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := utils.ParamsFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := utils.ParamsParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := FormParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := utils.FormString(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := utils.FormInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := utils.FormBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := utils.FormFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "file"}}
                        {{ParamInstance}}, err := utils.FormFile(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := utils.FormParser(webCtx, {{ParamInstance}})
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
                            {{if IsSingleton}}
                            ctx.{{ParamInstance}}(),
                            {{else if IsMultiple}}
                            ctx.New{{ParamInstance}}(),
                            {{end}}
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
                        {{ParamInstance}} := utils.Body(webCtx)
                    {{else if Type == "string"}}
                        {{ParamInstance}} := utils.BodyString(webCtx)
                    {{else}}
                        {{ParamInstance}} := &{{Package}}.{{ParamType}}{}
                        err := utils.BodyParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := HeaderParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := utils.Header(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := utils.HeaderInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := utils.HeaderBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := utils.HeaderFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := utils.HeaderParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := QueryParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := utils.Query(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := utils.QueryInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := utils.QueryBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := utils.QueryFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := utils.QueryParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := PathParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := utils.Params(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := utils.ParamsInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := utils.ParamsBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := utils.ParamsFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := utils.ParamsParser(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{end}}
                {{end}}
                
                {{range _, param := FormParams}}
                    {{if Type == "string"}}
                        {{ParamInstance}} := utils.FormString(webCtx, {{ParamInstance}})
                    {{else if Type == "int"}}
                        {{ParamInstance}}, err := utils.FormInt(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "bool"}}
                        {{ParamInstance}}, err := utils.FormBool(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "float64"}}
                        {{ParamInstance}}, err := utils.FormFloat(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else if Type == "file"}}
                        {{ParamInstance}}, err := utils.FormFile(webCtx, {{ParamInstance}})
                        if err != nil {
                            return err
                        }
                    {{else}}
                        {{ParamInstance}} := &model.Config{}
                        err := utils.FormParser(webCtx, {{ParamInstance}})
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
                            {{if IsSingleton}}
                            ctx.{{ParamInstance}}(),
                            {{else if IsMultiple}}
                            ctx.New{{ParamInstance}}(),
                            {{end}}
                        {{else}}
                            {{param.Instance}},
                        {{end}}
                    {{end}}
                )
            }
            {{end}}
        {{end}}
```
## 4. 使用规范

### 4.1. 包名规范
> 为了增加代码的可读性, 除了特殊类型的包, 如:`版本包`和`应用包`等之外, 其他情况下请务必保持包名与目录名一致。

下面是特殊类型包的命名规范:
* `版本包`：目录名格式为`v[\d.]+`, 字母v后跟数字; 它表示指定版本范围的包; <br/>`版本包`的包名只允许使用版本名和前一个的目录名, 如`github.com/gofiber/fiber/v2`, 只能使用`v2`或`fiber`
* `应用包`: 包名为`main`, 作为程序启动入口, 一般出现在mod模块的根目录. <br/>golang规定`main`不能被import导入, 所以main包即使定义了全局变量也无法被其他包访问;<br/>虽然`应用包`不能被import, 但是系统仍然会读取包内的注解.

### 4.2. 循环依赖问题

golang禁止两个包互相import导入, 为了避免它, 在设计上我们应当遵守<b>"声明与调用分离"</b>这一原则. 

具体做法如下:
* 应当准备两类包, 一类负责声明, 另一类负责调用; 调用包可以import依赖导入声明包, 但声明包禁止导入调用包;
* 声明包应当包含`@provide`,`@webAppProvide`,`@preConstruct`这些注解代码, 它们提供了实例的创建规则; 推荐包名为`model`; 
* 调用包应当包含`@postConstruct`,`@proxy`,`@middleware`,`@router`这些注解代码, 它们提供了依赖注入的函数回调; 推荐包名为`handler`;
