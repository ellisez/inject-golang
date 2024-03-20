package model

import (
	"go/ast"
	"go/token"
)

type Call struct {
	Params   []string
	Instance string
	Comment  string
}

type Field struct {
	Loc      token.Position
	Package  string
	Name     string   // 字段名
	Type     ast.Expr // 字段类型
	Instance string   // 实例名，默认同参数名
	Operator string   // 类型运算: '' | & | * | cast
	Index    int      // call索引
	Source   string   // 来源: '' | inject | func | call | query | path | header | body | formData | multipart
	Comment  string   // 原始注解
}

type Func struct {
	Loc      token.Position
	Package  string
	FuncName string
	Params   []*Field
	Recv     *Field
	Results  []*Field
	Calls    []*Call
}
