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

### 2.1. 不启用的注解
以`// !@`开头, 也就是注解前面多加个感叹号, 该注解不会被系统识别和解析
```go
// !@proxy
```

### 2.2. 代理函数注解

`代理函数注解`是让普通函数上支持依赖注入
```
// @proxy <代理方法名，默认同方法名>
// @override
// @import *<模块加载路径，必填>
// @injectParam *<参数名，必填> <实例名，默认同参数名> <运算, 默认""|&|*|cast>
// @injectRecv *<参数名，必填> <实例名，默认同参数名>
// @injectCtx *<参数名, 必填>
// @injectFunc *<参数名, 必填> <实例名，默认同参数名>
// @injectCall [*<参数名, 必填>, ...] <实例名，必填>
```

> `@proxy`让系统生成一个代理函数, 默认与原函数名同名. 代理函数可通过容器对象访问.
> 
> `@override`表示支持重载, 当遇到实例名相同时, 后者会覆盖前者; 默认是重载是关闭的, 同名时会报错.
>
> `@injectParam`用于参数的依赖注入;
>
> `@injectRecv`用于成员函数的结构体依赖注入;
>
> `@injectCtx`用于注入容器对象本身;
> 
> `@injectFunc`用于注入函数类型的参数;
> 
> `@injectCall`用于注入函数调用结果;
> 
> `@injectParam`支持类型转换, `&`表示取值的地址, `*`表示取地址对应的值, `cast`表示类型强转, 默认是""表示不进行转换.

> 注意: 未被依赖注入的参数则会保留到生成代理函数中;
>

### 2.3. 实例注解
`实例注解`指的是通过构造函数标记注解的方式来声明实例
```
// @provide <实例名，默认同返回类型> <singleton默认|multiple|argument> <实例type, 默认同返回类型>
// @override
// @order <创建顺序，数字或字符串>
// @import *<模块加载路径，必填>
// @injectParam *<参数名，必填> <实例名，默认同参数名> <指针运算, 默认""|&|*>
// @injectCtx *<参数名, 必填>
// @injectFunc *<参数名, 必填> <实例名，默认同参数名>
// @injectCall [*<参数名, 必填>, ...] <实例名，必填>
// @handler *<处理函数, 必填>
```

> 系统默认会加载容器包以及注解所在代码的包, 但如果有额外的包, 则需要使用`@import`声明.
> 
> `@provide`所标记的构造函数必须有且只有一个返回类型, 支持依赖注入, 但必须每个参数都要注入. 
> 
> `@provide`需要`@order`定义创建顺序, 以免还未完成初始化的实例被注入;
> 
> `@provide`模式有: `singleton`表示全局唯一, `multiple`表示可创建多个, `argument`表示仅在启动过程中存在;
> 
> `@override`表示支持重载, 当遇到实例名相同时, 后者会覆盖前者; 默认是重载是关闭的, 同名时会报错.
> 
> `@handler`所指向的函数要求必须无参数, 它会在实例创建后调用;
> 
> `@handler`可以携带包名, 例如: *model.Database, 它表示调用原始函数. 但如果不带包名, 则表示它是代理函数.
> 
> 一般我们不推荐直接使用原始函数, 而是使用它的代理函数, 这样能帮我们扩充其他依赖注入的参数; 
> 代理函数用法, 请参照[方法上注解](#23-方法上的注解-适用于所有方法)
>

### 2.4. WebApp注解 (提供了web服务器)
```
// @webProvide <instance，默认名为WebApp>
// @static *<访问路径，必填> *<匹配目录，必填> [特征: Compress|Download|Browse] <目录的Index文件> <过期时间MaxAge>
```

> `@webProvide`配置web应用实例, 不进行配置时系统默认会生成一个名为`WebApp`的实例.
> 
> 如果web应用实例未被代码中使用, 则系统不会生成WebApp实例, 这也是为了适配与非web项目.
> 
> web应用的启动函数, 格式为`instance + "Startup""`, 默认为`WebAppStartup`.
> 
> `@webProvide`所标记的原函数必须返回`host`,`port`,`err`三个参数.
> 
> `@static`用于配置静态资源文件, 如png,css,js,html等

### 2.5. 路由方法上的注解（参照swag注解）：
```
// @router *<Path必填> [Method: get|head|post|put|patch|delete|connect|options|trace]
// @webApp <WebApp，默认名为WebApp>
// @injectWebCtx *<参数名, 必填>
// @produce <返回格式: json | x-www-form-urlencoded | xml | plain | html | mpfd | json-api | json-stream | octet-stream | png | jpeg | gif>
// @param *<参数名，必填> *<取值类型，必填:query|path|header|body|formData> <接收类型> <必填与否> <参数说明>
```

> `@router`会让系统生成一个与函数同名的代理函数, 以完成参数的解析和注入, 也可以通过`@proxy`更改.
> 
> `@webApp`用于关联webApp实例, webApp由`@webProvide`提供, 默认实例名为"WebApp".
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

### 3.1. 创建实例
构建函数上的注解
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
> 构造函数, 要求必须只有一个返回类型
> 
> 范例中, 指定了实例类型`model.ServerInterface`, 它是个接口类型, 而真实创建类型为返回类型`*model.Server`.

原函数上的注解
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

> handler函数要求必须无参数, 上例中由于原函数完成特定的功能需要注入其他实例, 所以使用了代理函数.
> 

生成的代理函数如下:
```go
// Generate by annotations from handler.ServerAliasLoaded
func (ctx *Ctx) ServerAliasLoaded() {
    handler.ServerAliasLoaded(ctx, ctx.serverAlias.(*model.Server), ctx.database, &ctx.isReady, ctx.NewEvent(), ctx.NewListener())
}
```
> 生成的代理函数会保留那些未被注入的参数, 只有当所有参数都注入时, 才能生成一个无参的代理函数, 也就是`handler`的格式要求.
> 
> 参数`server`使用了`cast`运算, 因此生成代码中自动进行了强转. `cast`常用于接口与结构体的转换.
> 
> 参数`isReady`使用了`&`运算, 因此生成代码中自动进行了取地址操作. `isReady`是基本类型取指针是为了允许修改它.

### 3.2. web生成代码
注解配置
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
// @param username query string true 用户名
// @param password query string true 密码
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
生成代码
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
调用web的启动函数
```go
func main() {
    c := factory.New()
    err := c.WebAppStartup(3001)
    if err != nil {
        return
    }
}
```

> `@webProvide`用于定义web应用实例，不配置时系统默认会创建名为`WebApp`的实例，以及名为`WebAppStartup`的启动函数,
>
> 实例名通过`@webProvide`修改, 启动函数也会随着更改.
>
> `@router`与`@middleware`通过`@webApp`关联实例, 默认关联`WebApp`的实例名.
>
> `@webProvide`与`@webApp`实例名必须保持一致, 才能确保关联关系.

> <b>注意: `@middleware`与`@router`要求每个参数都必须都要配置依赖注入. 系统会自动创建出符合webApp调用格式的代码函数</b>
>
> 只有`@middleware`和`@router`可以注入`@webCtx`, 它表示本次请求上下文.
>
> webCtx更多使用技巧, 可阅读`*fiber.Ctx`相关文档
>

## 4. 推荐目录结构
inject-golang使用ctx简化了DDD(Domain Model Driven)目录结构.

```yaml
/main_domain                #当前域
# golang标准包
  /pkg                      #二进制包, .a/.so/.dll
  /bin                      #运行命令, .bat/.sh/.exe
  /vendor                   #外部依赖包，自动下再的，不应算到源代码中
  /third_party              #对第三方包的修改
# ddd分层
  # 用户接口层
  /ctx                      #代码生成，提供外部接口访问
  /interfaces               #提供外部数据结构
  # 应用层
  /application              #应用层, 外部接口实际入口, 仅组合现有服务, 不做具体实现.
    /startup                #外部可访问的启动项
    /controller             #外部可访问的请求项
    /service                #外部可访问的服务项
    ...
  # 当前领域模型
  /internal                 
    /entity                 #实体数据结构
    /vo                     #复合数据结构
    /service                #业务服务函数
    /repository             #数据存取函数
  /component                #基础设施层, 缓存/文件/网络通信等等
# 运行环境
  /conf                     #配置文件
  /web                      #web服务
    /handler                #手动路由处理函数, 不推荐
    /router                 #手动路由配置, 不推荐
    /middleware             #中间件函数, 拦截器/过滤器等
  /startup                  #编排启动项
    /config                 #配置启动项
    /db                     #数据库启动项
    /web                    #web服务启动项
    ...
  /rpc                      #rpc服务
  /utils                    #工具函数
  /mian.go                  #程序入口函数
/sub_domain                 #子领域, 与主域的目录结构一致
```

> 向外部提供访问有两种方式：服务方式和类库方式
> * 服务方式，就像微服务，需要启动和安装路由，然后通过发起请求调用服务；
> * 类库方式，则当作源代码引入，进行调用;
> 
> 一般两者只选其一, 不会同时出现.

> `/sub_domain`多领域模型是通过`go.work`来实现的, 每个子领域都和`/mian_domain`的目录结构相似.
> 
> 如果子领域是当作类库方式加载的, 那么`运行环境`相关的包是不需要的.

> 如果按照严格的ddd分层, `运行环境`包中`/handler`,`/router`,`/middleware`等提供调用的部分应属于`基础设施层`, 
> `/startup`等搭建运行环境的包应属于`用户接口层`. 
> 
> 但因为`运行环境`不具有复用度, 并且不必要增加目录层级, 所以放在了域的根目录上.
 
> 类库方式, 应当只`import`导入这三个包`/interfaces`,`/applcation`,`/ctx`, 其余包视为内部实现细节不应对外部提供.

> 虽然系统会自动生成外部访问, 但不可避免的需要提供注解和数据结构, 以告知系统生成代码的规则. 
> 接口中的数据结构, 应在`/interfaces`内提供, 注解应在`/application`内提供.
> 
> 按照ddd的思想, `/application`应该尽可能"薄", 所以它只是负责组织各个业务实现的调用, 而不会直接的实现业务代码.
> 

> `go:generate inject-golang`命令可以指定多个目录, 但最终生成的代码只会在当前目录的`/ctx`下.
> 这样做的目的是便于整合子包的注解, 同时注解也可以`@override`进行覆盖.

> 所以子领域模型类库方式下, 有两种实现方式: 一种是`单目录模式`, 各个子域自己生成`/ctx`包; 另一种方式是`多目录模式`, 所有子域都汇总到父域中生成`/ctx`包.
>
> 因为注解的代码会被生成的`/ctx`访问, 为了更好的隔离内部实现与外部访问, 做如下规定:
> * `单目录模式`下, 注解可直接标记到具体实现函数上, 外部只会访问`/ctx`包公开的接口, 内部实现不会访问到;
> * `多目录模式`下, 注解必须放在`/application`包, 并转接其他包具体实现代码, 不可在实现代码上加注解;

> DDD分层的复用性:  
> 1. 变更运行环境成本降低: 由于`/application`只是组装和转接没有具体实现,所以与`运行环境`无关, 当增加或者变更`运行环境`时, 只需要符合`/application`调用参数格式即可, 原来的代码依旧可以复用.
> <br>如: `/application`定义了`login()`登录函数需要参数用户名和密码, 虽然是web环境, 还是rpc服务或者其他服务, 解析参数用户名和密码各不相同, 但解析出来之后, 都会调用`login()`.
> 2. `Domain`领域层的复用: 由于领域层能独立完成业务功能的, 所以Domain层可以拷贝或引用到别的项目中, 以便组合出更加复杂的业务功能.

## 5. 使用规范

### 5.1. 包名规范
> 为了增加代码的可读性, 除了特殊类型的包, 如:`版本包`和`应用包`等之外, 其他情况下请务必保持包名与目录名一致。

下面是特殊类型包的命名规范:
* `版本包`：目录名格式为`v[\d.]+`, 字母v后跟数字; 它表示指定版本范围的包; <br/>`版本包`的包名应当是前一个的目录名, 如`github.com/gofiber/fiber/v2`中的`fiber`
* `应用包`: 包名为`main`, 作为程序启动入口, 一般出现在mod模块的根目录. <br/>golang规定`main`不能被import导入, 所以main包即使定义了全局变量也无法被其他包访问;<br/>虽然`应用包`不能被import, 但是系统仍然会读取包内的注解.

### 5.2. 循环依赖问题

golang禁止两个包互相import导入, 为了避免它, 在设计上我们应当遵守<b>"声明与调用分离"</b>这一原则. 

具体做法如下:
* 应当准备两类包, 一类负责声明, 另一类负责调用; 调用包可以import依赖导入声明包, 但声明包禁止导入调用包;
* 声明包应当包含`@provide`,`@webProvide`这些注解声明的结构体, 它们提供了实例的创建规则; 推荐包名为`model`; 
* 调用包应当包含`@handler`,`@proxy`,`@middleware`,`@router`这些注解代码, 它们提供了依赖注入的函数回调; 推荐包名为`handler`;

> 注意: 声明包是禁止使用`@injectCtx`注解的. 
> 
> 例如: 声明包的成员函数使用`@proxy`时, 不应注入`Ctx`上下文, 而是考虑替换成值注入和函数注入.
> 
> 值注入有`@injectParam`与`@injectRecv`等, 函数注入有`@injectFunc`.

### 5.3. 灵活的注解覆盖

在`go.work`跨模块工程里, 许多时候子模块定义的注解所生成的实例和代理函数, 不适配主模块, 传统的做法只能是新的实例和新的代理函数, 或者扩充子模块的代码.
新写代码则子模块的旧代码将无效造成冗余, 而扩充子模块又会造成子模块与主模块互相调用, 显然这两种做法都不太理想.
这是我们推荐使用`@override`注解覆盖.

注解覆盖可以做到同实例名, 后者覆盖前者, 而前者的注解不会在生成代码中留下任何痕迹.

例如, 子模块定义的结构体实例, 主模块中需要对结构体扩充一些字段, 这时只需要注解覆盖实例的类型即可.
