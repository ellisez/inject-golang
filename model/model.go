package model

import (
	"go/ast"
	"go/token"
)

type Module struct {
	Path    string            // the dir of go.mod
	Package string            // go.mod mod
	Version string            // go.mod version
	Require map[string]string // go.mod require
	Work    map[string]string // go.work
}

type Import struct {
	Name string
	Path string
}

type CommonFunc struct {
	Imports []*Import

	*Func
}

func NewCommonFunc() *CommonFunc {
	return &CommonFunc{
		Func: &Func{},
	}
}

type Comment struct {
	Comment string
	Args    []string
}

type Proxy struct {
	*CommonFunc
	Instance string
	Comment  string
}

func NewProxy() *Proxy {
	return &Proxy{CommonFunc: NewCommonFunc()}
}

type Method struct {
	From     *Func
	FuncName string
	Params   []*Field
	Results  []*Field
}

type Gen struct {
	Doc     []*ast.Comment
	Imports []*ast.ImportSpec
	Methods []*ast.FuncDecl
}
type Ctx struct {
	FileSet *token.FileSet

	SingletonInstances []Instance
	MultipleInstances  []Instance
	FuncInstances      []*Proxy
	MethodInstances    []*Proxy

	HasWebInstance bool
	*Gen
}

func NewCtx() *Ctx {
	return &Ctx{
		FileSet: token.NewFileSet(),
		Gen:     &Gen{},
	}
}

func (ctx *Ctx) MethodOf(funcName string) *ast.FuncDecl {
	for _, method := range ctx.Methods {
		if method.Name.String() == funcName {
			return method
		}
	}
	return nil
}

func (ctx *Ctx) SingletonOf(name string) Instance {
	for _, instance := range ctx.SingletonInstances {
		if instance.GetInstance() == name {
			return instance
		}
	}
	return nil
}
func (ctx *Ctx) MultipleOf(name string) Instance {
	for _, instance := range ctx.MultipleInstances {
		if instance.GetInstance() == name {
			return instance
		}
	}
	return nil
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
