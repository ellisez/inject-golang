package model

import (
	"go/ast"
	"go/token"
)

type Field struct {
	Name     string   // 字段名
	Type     ast.Expr // 字段类型
	Instance string   // 实例名，默认同参数名
	Pointer  string   // 指针运算: '' | & | *
	Source   string   // 来源: '' | inject | query | path | header | body | formData | multipart
	Comment  string   // 原始注解
}

type Func struct {
	Loc      token.Position
	Package  string
	FuncName string
	Params   []*Field
	Recv     *Field
	Results  []*Field
}
