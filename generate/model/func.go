package model

// FuncInfo
// @proxy <代理方法名，默认同方法名>
// @import <模块加载路径> <模块名>
// @injectParam <参数名> <实例名，默认同类名>
type FuncInfo struct {
	*PackageInfo

	Imports []*ImportInfo // import语句

	FuncName string // 函数名
	Proxy    string // 代理函数名，默认同函数名

	Recv *FieldInfo // 函数接收

	InjectParams []*FieldInfo // 注入字段
	NormalParams []*FieldInfo // 非注入字段
	Params       []*FieldInfo // 所有字段

	Results []*FieldInfo // 返回值

	ProxyComment string // @proxy注解
}
