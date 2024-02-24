### 2.1 结构体上的注解：
```
// @provide <实例名，默认同类名> <singleton默认|multiple>
// @constructor <构造函数名，默认New+类名>
```
### 2.2 结构体的属性注解：
```
// @inject <实例名，默认同类名>
```
### 2.3 方法上的注解:
```
// @proxy <代理方法名，默认同方法名>
```
### 2.4 路由方法上的注解（参照swag注解）：
```
// @route <path> <httpMethod: get|post>
// @param <参数名> <参数类型:query|path|header|body|formData|inject> <接收类型> <必填与否> <参数说明>
```
### 2.5 方法参数上的注解:
```
// @inject <实例名，默认同类名>
```

```
/provide
    |- __gen_container.go
        --------------------------------
        # gen segment: Struct #
        --------------------------------
        type ProvideContainer struct {
            {{range SingletonInstances}}
            Instance Name
            {{end}}
        }
      
        -----------------------------------
        # gen segment: Singleton instance #
        -----------------------------------
        func New() Container {
            container := &ProvideContainer{}
            {{range SingletonInstances}}
            container.{{Instance}} = &{{Name}}{}
            {{end}}
            
            {{range SingletonInstances}}
                {{range InjectFields}}
                container.{{Instance}}.{{FieldInstance}} = container.{{StructInstance}}
                {{end}}
            {{end}}
            return container
        }
    |- __gen_constructor.go
        ------------------------------------
        # gen segment: Multiple instance #
        ------------------------------------
        {{range MultipleInstances}}
            func (container *Container) {{Constructor}}(
                {{range NormalFields}}
                {{FieldInstance}} {{FieldType}},
                {{end}}
            ) *{{Type}} {
                {{Instance}} := &{{Type}}{}
                {{range Fields}}
                    {{if IsInject}}
                        {{Instance}}.{{FieldName}} = container.{{FieldInstance}}
                    {{else}}
                        {{Instance}}.{{FieldName}} = {{FieldInstance}}
                    {{end}}
                {{end}}
                return {{Instance}}
            }
        {{end}}
        
    |- __gen_func.go
        ------------------------------------
        # gen segment: Func inject #
        ------------------------------------
        {{range FuncInstances}}
            func (container *Container) {{Proxy}}(
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
                            container.{{ParamInstance}},
                        {{else}}
                            {{ParamInstance}},
                        {{end}}
                    {{end}}
                )
            }
        {{end}}
        
    |- __gen_method.go
        ------------------------------------
        # gen segment: Method inject #
        ------------------------------------
        func (container *Containe) {{Proxy}}(
            {{Recv.Name}} {{Recv.Type}},
            {{range NormalParams}}
            {{ParamInstance}} {{ParamType}},
            {{end}}
        ) (
            {{range Results}}
            {{ResultName}} {{ResultType}},
            {{end}}
        ) {
            return {{Package}}.{{FuncName}}(
                {{Recv.Name}},
                {{range Params}}
                    {{if IsInject}}
                        container.{{ParamInstance}},
                    {{else}}
                        {{ParamInstance}},
                    {{end}}
                {{end}}
            )
        }
```
