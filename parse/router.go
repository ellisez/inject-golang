package parse

import (
	"fmt"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"strings"
)

func (p *Parser) RouterParse(funcDecl *ast.FuncDecl, commonFunc *model.CommonFunc, comments []*model.Comment) {
	router := model.NewRouter()
	router.CommonFunc = commonFunc
	router.Instance = commonFunc.FuncName
	router.WebApp = "WebApp"

	commonFunc.Loc = p.Ctx.FileSet.Position(funcDecl.Pos())
	for _, comment := range comments {
		args := comment.Args
		argsLen := len(args)
		if argsLen == 0 {
			continue
		}
		annotateName := args[0]
		switch annotateName {
		case "@router":
			if argsLen < 2 {
				utils.Failuref("%s %s, Path must be specified", commonFunc.Loc.String(), comment.Text)
			}
			router.Path = args[1]

			if argsLen < 3 {
				utils.Failuref("%s %s, Methods must be specified", commonFunc.Loc.String(), comment.Text)
			}
			methods := args[2]
			if !strings.HasPrefix(methods, "[") || !strings.HasSuffix(methods, "]") {
				utils.Failuref("%s %s, Methods must be wrapped in []", commonFunc.Loc.String(), comment.Text)
			}
			methods = methods[1 : len(methods)-1]
			for _, method := range strings.Split(methods, ",") {
				router.Methods = append(router.Methods, utils.FirstToUpper(method))
			}

			router.Comment = comment.Text
		case "@param":
			webParamParse(router.WebParam, commonFunc, comment)
		case "@injectWebCtx":
			injectWebCtxParse(commonFunc, comment)
		case "@webApp":
			if argsLen >= 2 {
				router.WebApp = args[1]
			}
		}
	}

	for _, param := range router.Params {
		if param.Source == "" {
			utils.Failuref("%s %s, The ParamName \"%s\" does not allow non injection", commonFunc.Loc.String(), router.Comment, param.Name)
		}
	}

	instance, webApplication := p.Ctx.SingletonOf(router.WebApp)
	if instance != nil {
		if webApplication == nil {
			utils.Failuref(`%s %s, Conflict with "%s"`, commonFunc.Loc.String(), router.Comment, instance.Comment)
		}
		old, ok := webApplication.Routers[router.Instance]
		if ok {
			if !old.Override {
				utils.Failuref(`%s %s, Instance "%s" Duplicate declaration`, router.Loc.String(), router.Comment, router.Instance)
			}
			fmt.Printf(`Instance "%s" is Overrided by %s.%s`+"\n", router.WebApp, router.Package, router.FuncName)
		}
		webApplication.Routers[router.Instance] = router
	} else {
		newProvide := model.NewWebProvide()
		newProvide.Instance = router.WebApp

		newWebApplication := model.NewWebApplication()
		newWebApplication.Routers = map[string]*model.Router{
			router.Path: router,
		}

		p.Ctx.SingletonInstance.AddWeb(newProvide, newWebApplication)
	}

}
