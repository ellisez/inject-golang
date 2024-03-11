package model

import "strings"

// StructInfo
// @provide <Instance，default structName> <singleton default|multiple>
// @import *<Path, required> <Alias>
// @injectField *<FieldName, required> <Instance，default structName>
// @preConstruct *<before create call func, required>
// @postConstruct *<after created call func, required>
type StructInfo struct {
	// @inject <实例名，默认同类名>
	*PackageInfo
	Order                string
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

func (s *StructInfo) PrivateName() string {
	return strings.ToLower(s.Instance[0:1]) + s.Instance[1:]
}
func (s *StructInfo) Getter() string {
	switch s.Mode {
	case "singleton":
		return s.Instance
	case "multiple":
		return "New" + s.Instance
	}
	return ""
}

func (s *StructInfo) Setter() string {
	return "Set" + s.Instance
}

func NewStructInfo() *StructInfo {
	return &StructInfo{
		Mode: "singleton",
	}
}

type ImportInfo struct {
	Name string
	Path string
}
