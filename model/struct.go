package model

// StructInfo
// @provide <Instance，default structName> <singleton default|multiple>
// @import *<Path, required> <Alias>
// @injectField *<FieldName, required> <Instance，default structName>
// @preConstruct *<before create call func, required>
// @postConstruct *<after created call func, required>
type StructInfo struct {
	// @inject <实例名，默认同类名>
	*PackageInfo
	Imports              []*ImportInfo // import语句
	Name                 string        // 结构体名称
	Instance             string        // 实例名，默认同类名
	Mode                 string        // singleton默认|multiple
	PreConstruct         string        // 构造前调用函数
	PostConstruct        string        // 构造后调用函数
	Fields               []*FieldInfo  // 所有字段
	ProvideComment       string        // @provide注解
	PreConstructComment  string        // @preConstruct注解
	PostConstructComment string        // @postConstruct注解
}

func NewStructInfo() *StructInfo {
	return &StructInfo{}
}

type ImportInfo struct {
	Name string
	Path string
}
