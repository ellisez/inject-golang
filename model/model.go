package model

import (
	"go/ast"
)

type PackageInfo struct {
	Dirname string // 结构体目录
	Package string // 结构体包名
	Import  string // 用于import
}

type FieldInfo struct {
	Name     string   // 字段名
	Type     ast.Expr // 字段类型
	Instance string   // 实例名，默认同类名
	Source   string   // 来源: '' | inject | query | path | header | body | formData | multipart
	IsEmbed  bool     // 是否为内嵌属性，声明内嵌属性不具有字段名
	Comment  string   // 原始注解
}

type Mod struct {
	Path    string            // the dir of go.mod
	Package string            // go.mod mod
	Version string            // go.mod version
	Require map[string]string // go.mod require
	Work    map[string]string // go.work
}

type ModuleInfo struct {
	SingletonInstances []*StructInfo
	MultipleInstances  []*StructInfo
	FuncInstances      []*FuncInfo
	MethodInstances    []*FuncInfo

	WebAppInstances []*WebInfo
}

func NewModuleInfo() *ModuleInfo {
	return &ModuleInfo{
		SingletonInstances: make([]*StructInfo, 0),
		MultipleInstances:  make([]*StructInfo, 0),
		FuncInstances:      make([]*FuncInfo, 0),
		MethodInstances:    make([]*FuncInfo, 0),
		WebAppInstances:    make([]*WebInfo, 0),
	}
}

func (moduleInfo *ModuleInfo) HasFunc(funcName string) bool {
	for _, instance := range moduleInfo.MultipleInstances {
		if "New"+instance.Instance == funcName {
			return true
		}
	}

	for _, instance := range moduleInfo.FuncInstances {
		if instance.Proxy == funcName {
			return true
		}
	}

	for _, instance := range moduleInfo.MethodInstances {
		if instance.Proxy == funcName {
			return true
		}
	}
	return false
}

func (moduleInfo *ModuleInfo) GetSingleton(name string) *StructInfo {
	for _, instance := range moduleInfo.SingletonInstances {
		if instance.Instance == name {
			return instance
		}
	}
	return nil
}
func (moduleInfo *ModuleInfo) GetWebApp(name string) *WebInfo {
	for _, instance := range moduleInfo.WebAppInstances {
		if instance.WebApp == name {
			return instance
		}
	}
	return nil
}

func (moduleInfo *ModuleInfo) HasInstance(name string) bool {
	return moduleInfo.HasSingleton(name) || moduleInfo.HasWebApp(name)
}

func (moduleInfo *ModuleInfo) HasSingleton(name string) bool {
	return moduleInfo.GetSingleton(name) != nil
}

func (moduleInfo *ModuleInfo) HasWebApp(name string) bool {
	return moduleInfo.GetWebApp(name) != nil
}
