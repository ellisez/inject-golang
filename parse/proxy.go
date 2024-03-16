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
			proxy.Comment = comment.Comment
		}
	}

	if proxyOverride(p.Ctx, proxy) {
		return
	}
	if funcDecl.Recv == nil {
		p.Ctx.FuncInstances = append(p.Ctx.FuncInstances, proxy)
	} else {
		p.Ctx.MethodInstances = append(p.Ctx.MethodInstances, proxy)
	}
}

func proxyOverride(ctx *model.Ctx, proxy *model.Proxy) bool {
	for i, instance := range ctx.FuncInstances {
		if instance.Instance == proxy.Instance {
			if instance.Override {
				if proxy.Recv == nil {
					ctx.FuncInstances[i] = proxy
				} else {
					ctx.MethodInstances = append(ctx.MethodInstances, proxy)
				}
				fmt.Printf(`Proxy "%s" is Overrided by %s.%s`+"\n", proxy.Instance, proxy.Package, proxy.FuncName)
				return true
			} else {
				utils.Failuref(`%s %s, Proxy "%s" Duplicate declaration`, proxy.Loc.String(), proxy.Comment, proxy.Instance)
			}
		}
	}
	for i, instance := range ctx.MethodInstances {
		if instance.Instance == proxy.Instance {
			if instance.Override {
				if proxy.Recv == nil {
					ctx.FuncInstances = append(ctx.FuncInstances, proxy)
				} else {
					ctx.MethodInstances[i] = proxy
				}
				fmt.Printf(`Instance "%s" is Overrided by %s.%s`+"\n", proxy.Instance, proxy.Package, proxy.FuncName)
				return true
			} else {
				utils.Failuref(`%s %s, Instance "%s" Duplicate declaration`, proxy.Loc.String(), proxy.Comment, proxy.Instance)
			}
		}
	}
	return false
}
