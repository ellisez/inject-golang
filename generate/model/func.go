package model

type FuncInfo struct {
	*PackageInfo

	FuncName string // 函数名
	Proxy    string // 代理函数名，默认同函数名

	Recv *ParamInfo // 函数接收

	InjectParams []*ParamInfo // 注入字段
	NormalParams []*ParamInfo // 非注入字段
	Params       []*ParamInfo // 所有字段

	Results []*ParamInfo // 返回值
}
