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
	IsInject bool     // 是否注入
	IsEmbed  bool     // 是否为内嵌属性，声明内嵌属性不具有字段名
	Comment  string   // 原始注解
}

type ModuleInfo struct {
	Dirname            string // 根目录, 用于写文件路径
	Mod                string // golang模块名, 用于import
	SingletonInstances []*StructInfo
	MultipleInstances  []*StructInfo
	FuncInstances      []*FuncInfo
	MethodInstances    []*FuncInfo
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

func (moduleInfo *ModuleInfo) HasStruct(name string) bool {
	return moduleInfo.GetSingleton(name) != nil
}
