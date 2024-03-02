package model

// FuncInfo (use for all func)
// @proxy <Instance，default funcName>
// @import *<Path, required> <Alias>
// @injectParam *<ParamName, required> <Instance，default paramName>
// @injectRecv *<ParamName, required> <Instance，default paramName>
type FuncInfo struct {
	*PackageInfo

	Imports []*ImportInfo // import语句

	FuncName string // 函数名
	Proxy    string // 代理函数名，默认同函数名

	Recv *FieldInfo // 函数接收

	Params []*FieldInfo // 所有字段

	Results []*FieldInfo // 返回值

	ProxyComment string // @proxy注解
}

func NewFuncInfo() *FuncInfo {
	return &FuncInfo{
		Imports: make([]*ImportInfo, 0),
		Params:  make([]*FieldInfo, 0),
		Results: make([]*FieldInfo, 0),
	}
}
