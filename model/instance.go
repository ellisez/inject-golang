package model

import "go/ast"

type Provide struct {
	*CommonFunc
	Mode string

	Instance    string
	Constructor ast.Expr
	Type        ast.Expr
	Handler     string
}

func NewProvide() *Provide {
	return &Provide{
		Mode:       "singleton",
		CommonFunc: NewCommonFunc()}
}
