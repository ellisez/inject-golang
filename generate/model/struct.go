package model

// StructInfo
// @provide <实例名，默认同类名> <singleton默认|multiple>
// @constructor <构造函数名，默认New+类名>
type StructInfo struct {
	*PackageInfo
	Name         string       // 结构体名称
	Instance     string       // 实例名，默认同类名
	Mode         string       // singleton默认|multiple
	Constructor  string       // 构造函数名，默认New+类名
	InjectFields []*FieldInfo // 注入字段
	NormalFields []*FieldInfo // 非注入字段
	Fields       []*FieldInfo // 所有字段
}
