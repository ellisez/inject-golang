package parse

import (
	"github.com/ellisez/inject-golang/model"
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

	if funcDecl.Recv == nil {
		p.Ctx.FuncInstances = append(p.Ctx.FuncInstances, proxy)
	} else {
		p.Ctx.MethodInstances = append(p.Ctx.MethodInstances, proxy)
	}
}
