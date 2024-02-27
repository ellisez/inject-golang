# inject-golang
Provide container of DI(Dependency Injection) for golang.


### 2.1 结构体上的注解：
```
// @provide <实例名，默认同类名> <singleton默认|multiple>
// @import <模块加载路径> <模块名>
// @injectField <字段名> <实例名，默认同类名>
// @preConstruct <构造前调用函数>
// @postConstruct <构造后调用函数>
```

### 2.2 结构体的属性注解：
```
// @inject <实例名，默认同类名>
```
### 2.3 方法上的注解:
```
// @proxy <代理方法名，默认同方法名>
// @import <模块加载路径> <模块名>
// @injectParam <参数名> <实例名，默认同类名>
```
### 2.4 路由方法上的注解（参照swag注解）：
```
// @route <path> <httpMethod: get|post>
// @import <模块加载路径> <模块名>
// @paramInject <参数名> <实例名，默认同类名>
// @param <参数名> <参数类型:query|path|header|body|formData> <接收类型> <必填与否> <参数说明>
```

### 2.5 方法参数上的注解:
```
// @inject <实例名，默认同类名>
```
`@inject @injectParam @injectField 预留实例名"ctx"时表示注入上下文容器本身`
```
/provide
    |- gen_ctx.go
        --------------------------------
        # gen segment: Struct #
        --------------------------------
        type Ctx struct {
            {{range SingletonInstances}}
            Instance Name
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
        func (ctx *Containe) {{Proxy}}(
            {{Recv.Name}} {{Recv.Type}},
            {{range NormalParams}}
            {{ParamInstance}} {{ParamType}},
            {{end}}
        ) (
            {{range Results}}
            {{ResultName}} {{ResultType}},
            {{end}}
        ) {
            return {{Recv.Name}}.{{FuncName}}(
                {{Recv.Name}},
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
```
