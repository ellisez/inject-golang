package model

// StructInfo
// @provide <实例名，默认同类名> <singleton默认|multiple>
// @import <模块加载路径，必填> <模块名>
// @injectField <字段名, 必填> <实例名，默认同类名>
// @preConstruct <构造前调用函数，必填>
// @postConstruct <构造后调用函数，必填>
type StructInfo struct {
	// @inject <实例名，默认同类名>
	*PackageInfo
	Imports              []*ImportInfo // import语句
	Name                 string        // 结构体名称
	Instance             string        // 实例名，默认同类名
	Mode                 string        // singleton默认|multiple
	PreConstruct         string        // 构造前调用函数
	PostConstruct        string        // 构造后调用函数
	InjectFields         []*FieldInfo  // 注入字段
	NormalFields         []*FieldInfo  // 非注入字段
	Fields               []*FieldInfo  // 所有字段
	ProvideComment       string        // @provide注解
	PreConstructComment  string        // @preConstruct注解
	PostConstructComment string        // @postConstruct注解
}

type ImportInfo struct {
	Name string
	Path string
}
