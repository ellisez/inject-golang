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
func (provide *Provide) GetConstructor() ast.Expr {
	return provide.Constructor
}
func (provide *Provide) GetOverride() bool {
	return provide.Override
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
	GetConstructor() ast.Expr
	GetOverride() bool
}
