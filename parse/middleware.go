package parse

import (
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

func (p *Parser) MiddlewareParse(funcDecl *ast.FuncDecl, commonFunc *model.CommonFunc, comments []*model.Comment) {
	middleware := model.NewMiddleware()
	middleware.CommonFunc = commonFunc
	middleware.Instance = commonFunc.FuncName
	middleware.WebApp = "WebApp"

	commonFunc.Loc = p.Ctx.FileSet.Position(funcDecl.Pos())
	for _, comment := range comments {
		args := comment.Args
		argsLen := len(args)
		if argsLen == 0 {
			continue
		}
		annotateName := args[0]
		switch annotateName {
		case "@middleware":
			if argsLen < 2 {
				utils.Failuref("%s %s, Path must be specified", commonFunc.Loc.String(), comment.Text)
			}
			middleware.Path = args[1]

			middleware.Comment = comment.Text
		case "@param":
			webParamParse(middleware.WebParam, commonFunc, comment)
		case "@injectWebCtx":
			injectWebCtxParse(commonFunc, comment)
		case "@webApp":
			if argsLen >= 2 {
				middleware.WebApp = args[1]
			}
		}
	}

	for _, param := range middleware.Params {
		if param.Source == "" {
			utils.Failuref("%s %s, The ParamName \"%s\" does not allow non injection", commonFunc.Loc.String(), middleware.Comment, param.Name)
		}
	}

	instance := p.Ctx.SingletonOf(middleware.WebApp)
	if instance != nil {
		webInstance, ok := instance.(*model.WebInstance)
		if !ok {
			utils.Failuref(`%s %s, Conflict with "%s"`, commonFunc.Loc.String(), middleware.Comment, instance.GetComment())
		}
		old, ok := webInstance.Middlewares[middleware.Instance]
		if ok && !old.Override {
			utils.Failuref(`%s %s, Instance "%s" Duplicate declaration`, middleware.Loc.String(), middleware.Comment, middleware.Instance)
		}
		webInstance.Middlewares[middleware.Instance] = middleware
	} else {
		webInstance := model.NewWebInstance()
		webInstance.Middlewares = map[string]*model.Middleware{
			middleware.Path: middleware,
		}
		webInstance.Instance = middleware.WebApp

		p.Ctx.SingletonInstances.Add(webInstance)
	}
	p.Ctx.HasWebInstance = true
}
