package model

import (
	"go/ast"
	"go/token"
)

type Gen struct {
	Doc          []*ast.Comment
	Imports      []*ast.ImportSpec
	Methods      map[string]*ast.FuncDecl
	InjectCtxMap map[string][]*Field
}
type Ctx struct {
	FileSet        *token.FileSet
	PackageMapping map[string]string

	SingletonInstances *InstanceMap

	MultipleInstances *InstanceMap

	FuncInstances *ProxyMap

	MethodInstances *ProxyMap

	HasWebInstance bool
	*Gen
}

func NewCtx() *Ctx {
	return &Ctx{
		FileSet:            token.NewFileSet(),
		PackageMapping:     map[string]string{},
		SingletonInstances: NewInstanceMap(),
		MultipleInstances:  NewInstanceMap(),
		FuncInstances:      NewProxyMap(),
		MethodInstances:    NewProxyMap(),
		Gen: &Gen{
			Methods:      map[string]*ast.FuncDecl{},
			InjectCtxMap: map[string][]*Field{},
		},
	}
}

func (ctx *Ctx) MethodOf(funcName string) *ast.FuncDecl {
	return ctx.Methods[funcName]
}

func (ctx *Ctx) SingletonOf(name string) Instance {
	return ctx.SingletonInstances.Get(name)
}
func (ctx *Ctx) MultipleOf(name string) Instance {
	return ctx.MultipleInstances.Get(name)
}

func (ctx *Ctx) InstanceOf(name string) Instance {
	instance := ctx.SingletonOf(name)
	if instance != nil {
		return instance
	}
	instance = ctx.MultipleOf(name)
	if instance != nil {
		return instance
	}
	return nil
}
