# inject-golang
提供golang版依赖注入容器

语言: [English](README.md) 中文

### 1. 安装与运行

```shell
go install github.com/ellisez/inject-golang
```
### 1.1. 配置生成器
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
> //go:generate inject-golang

### 1.2. 运行生成器
```shell
go generate -run inject-golang
```
## 2. 注解
### 2.1. 结构体上的注解：
```
// @provide <实例名，默认同类名> <singleton默认|multiple>
// @import *<模块加载路径，必填> <模块名>
// @injectField *<字段名，必填> <实例名，默认同类名>
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
> 但一般我们不推荐直接使用唯一参数的函数, 而是使用代理函数, 帮我们扩充依赖注入的参数; 
> 代理函数用法, 请参照[方法上注解](#23-方法上的注解-适用于所有方法)
>
>

### 2.2. 结构体的属性注解：
```
// @inject <实例名，默认同类名>
```
### 2.3. 方法上的注解 (适用于所有方法)
```
// @proxy <代理方法名，默认同方法名>
// @import *<模块加载路径，必填> <模块名>
// @injectParam *<参数名，必填> <实例名，默认同类名>
// @injectRecv *<参数名，必填> <实例名，默认同类名>
// @injectCtx *<参数名, 必填>
```

> `@proxy`注解用于标记该函数支持依赖注入, 最终系统会自动生成一个与函数名同名代理函数, 可通过容器对象获得代理函数.
> 
> `@injectParam`用于参数的依赖注入;
> 
> `@injectRecv`用于从属结构体类型的依赖注入;
> 
> `@injectCtx`把容器本身注入到指定参数中;

### 2.4. WebApp注解 (提供了web服务器)
```
// @webAppProvide <WebApp，默认名为WebApp>
// @static *<访问路径，必填> *<匹配目录，必填> [特征: Compress|Download|Browse] <目录的Index文件> <过期时间MaxAge>
```

> `@webAppProvide`配置webApp, 不进行配置时系统默认会生成一个名为WebApp的实例. 
> 
> 如果webApp未被代码中使用, 则系统不会生成WebApp实例, 这也是为了使用与非web项目.
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

> `@router`系统会生成一个与函数同名的代理函数, 以完成参数的解析和注入.
> 
> `@webApp`用于指明路由注册在哪个webApp实例上, webApp由`@webAppProvide`进行配置, 默认会有一个名为WebApp实例, 当然名字也是可以更改的.
> 
> `@injectWebCtx`用于指定某个参数来接收当前请求的webCtx.
> 
> `@injectWebCtx`只有`@middleware`和`@router`才可以使用, 因为只有它们能够处理请求.
> 
> `@produce`用于定义返回数据类型
> 
> `@param`用于路由参数的解析, 支持下列格式:
> * query: Get参数, 如/index.html?a=1;
> * path: 路由参数, 如/article/:id;
> * header: 头部信息参数;
> * body: 请求body二进制流, 注意只能由一个body参数;
> * formData: multipart/form方式提交的数据;

### 2.6. 路由中间件上的注解:
```
// @middleware *<Path必填>
// @webApp <WebApp，默认名为WebApp>
// @injectWebCtx *<参数名, 必填>
// @param *<参数名，必填> *<参数类型，必填:query|path|header|body|formData> <接收类型> <必填与否> <参数说明>
```

> `@middleware`系统会生成一个与函数同名的代理函数, 以完成参数的解析和注入.
>

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
preConstruct函数, 要求必须是无参并且返回类型必须和结构体类型一样
```go
func PrepareWebCtxAlias() *WebApp {
    return &WebApp{}
}
```
postConstruct函数, 要求参数必须是结构体类型
```go
func WebCtxAliasLoaded(webApp *WebApp) {
	
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

> 由于postConstruct函数必须接收一个与注解的结构体一样的参数，而proxy代理函数对未注入的参数会保留到生成的函数中，
> 
> 所以我们推荐，postConstruct直接proxy生成的函数，而不是直接指向原函数。
> 利用原函数保留未注入结构体参数，来促成生成的函数满足postConstruct的要求。

### 3.3. 目录结构
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
                    ctx.{{Instance}}.{{FieldInstance}} = ctx.{{StructInstance}}
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
            func (ctx *Ctx) New{{Instance}}(
                {{range NormalFields}}
                {{FieldInstance}} {{FieldType}},
                {{end}}
            ) *{{Type}} {
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
                        {{Instance}}.{{FieldName}} = ctx.{{FieldInstance}}
                        {{end}}
                    {{else}}
                        {{Instance}}.{{FieldName}} = {{FieldInstance}}
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
                            ctx.{{ParamInstance}},
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
                        ctx.{{ParamInstance}},
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
