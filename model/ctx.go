package model

import (
	"go/ast"
	"go/token"
)

type Extra struct {
	Doc     []*ast.Comment
	Imports []*ast.ImportSpec
	Methods map[string]*ast.FuncDecl
}

type Key struct {
	Instance string
	Order    string
	Type     string // singleton, multiple, webApplication, argument, func, method
}
type Ctx struct {
	FileSet *token.FileSet

	SingletonInstance *SingletonInstance

	MultipleInstance *MultipleInstance

	FuncInstance *FuncInstance

	// other
	*Extra
}

func NewCtx() *Ctx {
	return &Ctx{
		FileSet:           token.NewFileSet(),
		SingletonInstance: newCtxSingletonInstance(),
		MultipleInstance:  newCtxMultipleInstance(),
		FuncInstance:      newCtxFuncInstance(),
		Extra: &Extra{
			Methods: map[string]*ast.FuncDecl{},
		},
	}
}

func (ctx *Ctx) MethodOf(funcName string) *ast.FuncDecl {
	return ctx.Methods[funcName]
}

func (ctx *Ctx) SingletonOf(name string) (*Provide, *WebApplication) {
	return ctx.SingletonInstance.Get(name)
}

func (ctx *Ctx) MultipleOf(name string) *Provide {
	return ctx.MultipleInstance.Get(name)
}

func (ctx *Ctx) InstanceOf(name string) (*Provide, *WebApplication) {
	instance, w := ctx.SingletonOf(name)
	if instance != nil {
		return instance, w
	}
	instance = ctx.MultipleOf(name)
	if instance != nil {
		return instance, nil
	}
	return nil, nil
}

func (ctx *Ctx) FuncOf(name string) *Proxy {
	return ctx.FuncInstance.Get(name)
}
