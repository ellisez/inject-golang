package model

type PackageInfo struct {
	Dirname string // 结构体目录
	Package string // 结构体包名
}

type FieldInfo struct {
	Name     string // 字段名
	Type     string // 字段类型
	Instance string // 实例名，默认同类名
	IsInject bool   // 是否注入
}

type ParamInfo struct {
	Name     string // 参数名
	Type     string // 参数类型
	Instance string // 实例名，默认同类名
	IsInject bool   // 是否注入
}
type AnnotateInfo struct {
	SingletonInstances []*StructInfo
	MultipleInstances  []*StructInfo
	FuncInstances      []*FuncInfo
	MethodInstances    []*FuncInfo
}
