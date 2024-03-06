package model

type RouterParam struct {
	QueryParams  []*FieldInfo // query取值
	PathParams   []*FieldInfo // path取值
	HeaderParams []*FieldInfo // header取值
	BodyParam    *FieldInfo   // body取值
	FormParams   []*FieldInfo // formData取值
}

func NewRouterParam() *RouterParam {
	return &RouterParam{
		QueryParams:  make([]*FieldInfo, 0),
		PathParams:   make([]*FieldInfo, 0),
		HeaderParams: make([]*FieldInfo, 0),
		FormParams:   make([]*FieldInfo, 0),
	}
}

// MiddlewareInfo
// @webApp <WebApp，default WebApp>
// @static *<Path, required> *<Path, required> [Features: Compress|Download|Browse] <Index> <MaxAge>
type MiddlewareInfo struct {
	*FuncInfo

	*RouterParam

	WebApp            string // 所属WebApp，默认WebApp
	Path              string // 路径
	Handle            string // 处理函数名fiber.Handler，当前函数名
	MiddlewareComment string
}

func NewMiddlewareInfoFromFuncInfo(funcInfo *FuncInfo) *MiddlewareInfo {
	return &MiddlewareInfo{
		FuncInfo:    funcInfo,
		RouterParam: NewRouterParam(),
		WebApp:      "WebApp",
	}
}

// RouterInfo
// @webApp <WebApp，default WebApp>
// @route *<Path, required> [Method: get|post]
// @param *<ParamName, required> <type:query|path|header|body|formData> <DataType> <IsRequired> <Description>
type RouterInfo struct {
	*FuncInfo

	*RouterParam

	WebApp string // 所属WebApp，默认WebApp

	Methods       []string // 请求方式get|post|put|patch
	Path          string   // 路径
	RouterComment string   // @router注解
}

func NewRouterInfoFromFuncInfo(funInfo *FuncInfo) *RouterInfo {
	return &RouterInfo{
		FuncInfo:    funInfo,
		RouterParam: NewRouterParam(),
		WebApp:      "WebApp",
		Methods:     make([]string, 0),
	}
}

type StaticResource struct {
	Path          string
	Dirname       string
	Features      []string
	Index         string
	MaxAge        int
	StaticComment string
}

func NewStaticResource() *StaticResource {
	return &StaticResource{
		Features: make([]string, 0),
	}
}

// WebInfo
// @webAppProvide <实例名，默认webApp>
// @static <Path> <Path> [Features: Compress|Download|Browse] <Index> <MaxAge>
type WebInfo struct {
	*FuncInfo

	WebApp        string            // WebApp实例名
	Statics       []*StaticResource // 静态资源
	Middlewares   []*MiddlewareInfo // 组内中间件
	Routers       []*RouterInfo     // 组内路由
	WebAppComment string            // @webApp注解
}

func NewWebInfo() *WebInfo {
	funInfo := NewFuncInfo()
	funInfo.Proxy = "WebAppStartup"
	return &WebInfo{
		FuncInfo:    funInfo,
		WebApp:      "WebApp",
		Statics:     make([]*StaticResource, 0),
		Middlewares: make([]*MiddlewareInfo, 0),
		Routers:     make([]*RouterInfo, 0),
	}
}

func NewWebInfoFromFunc(funInfo *FuncInfo) *WebInfo {
	funInfo.Proxy = "WebAppStartup"
	return &WebInfo{
		FuncInfo:    funInfo,
		WebApp:      "WebApp",
		Statics:     make([]*StaticResource, 0),
		Middlewares: make([]*MiddlewareInfo, 0),
		Routers:     make([]*RouterInfo, 0),
	}
}
