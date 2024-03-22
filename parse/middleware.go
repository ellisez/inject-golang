package parse

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
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

	instance, webApplication := p.Ctx.SingletonOf(middleware.WebApp)
	if instance != nil {
		if webApplication == nil {
			utils.Failuref(`%s %s, Conflict with "%s"`, commonFunc.Loc.String(), middleware.Comment, instance.Comment)
		}
		old, ok := webApplication.Middlewares[middleware.Instance]
		if ok {
			if !old.Override {
				utils.Failuref(`%s %s, Instance "%s" Duplicate declaration`, middleware.Loc.String(), middleware.Comment, middleware.Instance)
			}
			fmt.Printf(`Middleware "%s" is Overrided by %s.%s`+"\n", middleware.Instance, ImportAliasMap[middleware.Package].Path, middleware.FuncName)
		}
		webApplication.Middlewares[middleware.Instance] = middleware
	} else {
		newProvide := model.NewWebProvide()
		newProvide.Instance = middleware.WebApp

		newWebApplication := model.NewWebApplication()
		newWebApplication.Middlewares = map[string]*model.Middleware{
			middleware.Path: middleware,
		}

		p.Ctx.SingletonInstance.AddWeb(newProvide, newWebApplication)
	}
}
