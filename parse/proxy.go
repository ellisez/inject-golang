package parse

import (
	"fmt"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

func (p *Parser) ProxyParse(funcDecl *ast.FuncDecl, commonFunc *model.CommonFunc, comments []*model.Comment) {
	proxy := model.NewProxy()
	proxy.CommonFunc = commonFunc
	proxy.Instance = commonFunc.FuncName
	for _, comment := range comments {
		args := comment.Args
		argsLen := len(args)
		if argsLen == 0 {
			continue
		}
		annotateName := args[0]
		switch annotateName {
		case "@proxy":
			if argsLen >= 2 {
				instance := args[1]
				if instance != "" && instance != "_" {
					proxy.Instance = instance
				}
			}
			proxy.Comment = comment.Text
		}
	}

	addProxy(p.Ctx, proxy)
}

func addProxy(ctx *model.Ctx, proxy *model.Proxy) {
	instance := ctx.FuncOf(proxy.Instance)
	if instance != nil {
		if !instance.Override {
			utils.Failuref(`%s %s, Proxy "%s" Duplicate declaration`, proxy.Loc.String(), proxy.Comment, proxy.Instance)
		}
		fmt.Printf(`Proxy "%s" is Overrided by %s.%s`+"\n", proxy.Instance, proxy.Package, proxy.FuncName)
		ctx.FuncInstance.Replace(proxy)
	} else {
		ctx.FuncInstance.Add(proxy)
	}
}
