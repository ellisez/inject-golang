package model

import "go/ast"

type Provide struct {
	*CommonFunc
	Mode string

	Order    string
	Instance string
	Type     ast.Expr
	Handler  string
	Comment  string
}

func NewProvide() *Provide {
	return &Provide{CommonFunc: NewCommonFunc()}
}

func (provide *Provide) GetOrder() string {
	return provide.Order
}
func (provide *Provide) GetInstance() string {
	return provide.Instance
}
func (provide *Provide) GetComment() string {
	return provide.Comment
}
func (provide *Provide) GetImports() []*Import {
	return provide.Imports
}

func (provide *Provide) GetFunc() *Func {
	return provide.Func
}
func (provide *Provide) GetMode() string {
	return provide.Mode
}

func (provide *Provide) GetType() ast.Expr {
	return provide.Type
}

func (provide *Provide) GetHandler() string {
	return provide.Handler
}

type Instance interface {
	GetOrder() string

	GetInstance() string
	GetComment() string
	GetImports() []*Import
	GetFunc() *Func
	GetType() ast.Expr
	GetMode() string
	GetHandler() string
}
