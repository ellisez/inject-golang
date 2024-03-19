package parse

import (
	"fmt"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"strconv"
	"strings"
)

func (p *Parser) WebParse(funcDecl *ast.FuncDecl, commonFunc *model.CommonFunc, comments []*model.Comment) {
	provide := model.NewWebProvide()
	provide.CommonFunc = commonFunc

	webApp := model.NewWebApplication()

	commonFunc.Loc = p.Ctx.FileSet.Position(funcDecl.Pos())

	webValidate(provide)
	for _, comment := range comments {
		args := comment.Args
		argsLen := len(args)
		if argsLen == 0 {
			continue
		}
		annotateName := args[0]
		switch annotateName {
		case "@webProvide":
			if argsLen >= 2 {
				instance := args[1]
				if instance != "" && instance != "_" {
					if utils.IsFirstLower(instance) {
						utils.Failuref(`%s %s, Instance "%s" must be capitalized with the first letter`, commonFunc.Loc.String(), provide.Comment, instance)
					}
					provide.Instance = instance
				}
			}
			provide.Comment = comment.Text
		case "@static":
			if argsLen < 2 {
				utils.Failuref("%s %s, Path must be specified", commonFunc.Loc.String(), comment.Text)
			}
			path := args[1]
			if argsLen < 3 {
				utils.Failuref("%s %s, Dirname must be specified", commonFunc.Loc.String(), comment.Text)
			}
			dirname := args[2]

			resource := &model.WebResource{
				Path:    path,
				Dirname: dirname,
			}
			webApp.Resources[resource.Path] = resource

			if argsLen >= 4 {
				features := args[3]
				if !strings.HasPrefix(features, "[") || !strings.HasSuffix(features, "]") {
					utils.Failuref("%s %s, Features must be wrapped in []", commonFunc.Loc.String(), comment.Text)
				}
				features = features[1 : len(features)-1]
				splitStr := strings.Split(features, ",")
				for _, feature := range splitStr {
					resource.Features = append(resource.Features, strings.TrimSpace(feature))
				}
			}
			if argsLen >= 5 {
				resource.Index = args[4]
			}
			if argsLen >= 6 {
				maxAge, err := strconv.Atoi(args[5])
				if err != nil {
					utils.Failuref("%s %s, MaxAge must be a number", commonFunc.Loc.String(), comment.Text)
				}
				resource.MaxAge = maxAge
			}
			resource.Comment = comment.Text
		}
	}

	instance, webApplication := p.Ctx.SingletonOf(provide.Instance)
	if instance != nil {
		if webApplication == nil {
			utils.Failuref(`%s %s, Conflict with "%s"`, commonFunc.Loc.String(), instance.Comment, instance.Comment)
		}
		if instance.FuncName != "" {
			if !instance.Override {
				utils.Failuref(`%s %s, Instance "%s" Duplicate declaration`, provide.Loc.String(), provide.Comment, provide.Instance)
			}
			fmt.Printf(`Instance "%s" is Overrided by %s.%s`+"\n", provide.Instance, provide.Package, provide.FuncName)
		}
		instance.Comment = provide.Comment
		instance.Imports = append(instance.Imports, provide.Imports...)
		instance.Func = provide.Func
		webApplication.Resources = webApp.Resources
	} else {
		p.Ctx.SingletonInstance.AddWeb(provide, webApp)
	}
}

func webValidate(provide *model.Provide) {
	if len(provide.Results) != 3 {
		utils.Failuref(`%s %s, Illegal webProvide function, returns [host, port, err]`, provide.Loc.String(), provide.Comment)
	}
	host, ok := provide.Results[0].Type.(*ast.Ident)
	if !ok {
		utils.Failuref(`%s %s, Illegal webProvide function, 1st returns is not "string"`, provide.Loc.String(), provide.Comment)
	}
	if host.String() != "string" {
		utils.Failuref(`%s %s, Illegal webProvide function, 1st returns is not "string"`, provide.Loc.String(), provide.Comment)
	}

	port, ok := provide.Results[1].Type.(*ast.Ident)
	if !ok {
		utils.Failuref(`%s %s, Illegal webProvide function, 2st returns is not "uint"`, provide.Loc.String(), provide.Comment)
	}
	if port.String() != "uint" {
		utils.Failuref(`%s %s, Illegal webProvide function, 2st returns is not "uint"`, provide.Loc.String(), provide.Comment)
	}

	err, ok := provide.Results[2].Type.(*ast.Ident)
	if !ok {
		utils.Failuref(`%s %s, Illegal webProvide function, 3st returns is not "error"`, provide.Loc.String(), provide.Comment)
	}
	if err.String() != "error" {
		utils.Failuref(`%s %s, Illegal webProvide function, 3st returns is not "error"`, provide.Loc.String(), provide.Comment)
	}

}

func webParamParse(webParam *model.WebParam, commonFunc *model.CommonFunc, comment *model.Comment) {
	args := comment.Args
	argsLen := len(args)
	if argsLen < 2 {
		utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Text)
	}
	paramName := args[1]
	paramInfo := utils.FindParam(commonFunc.Func, paramName)
	if paramInfo == nil {
		utils.Failuref("%s %s, ParamName not found", commonFunc.Loc.String(), comment.Text)
	}

	if argsLen < 3 {
		utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Text)
	}
	paramSource := args[2]
	paramInfo.Comment = comment.Text
	paramInfo.Source = paramSource
	addWebParam(webParam, paramInfo, commonFunc)
}

func injectWebCtxParse(commonFunc *model.CommonFunc, comment *model.Comment) {
	args := comment.Args
	argsLen := len(args)
	if argsLen < 2 {
		utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Text)
	}
	paramName := args[1]
	paramInfo := utils.FindParam(commonFunc.Func, paramName)
	if paramInfo == nil {
		utils.Failuref("%s %s, ParamName not found", commonFunc.Loc.String(), comment.Text)
	}

	paramInfo.Comment = comment.Text
	paramInfo.Source = "webCtx"
}

func addWebParam(webParam *model.WebParam, param *model.Field, commonFunc *model.CommonFunc) {
	switch param.Source {
	case "query":
		webParam.QueryParams = append(webParam.QueryParams, param)
	case "path":
		webParam.PathParams = append(webParam.PathParams, param)
	case "header":
		webParam.HeaderParams = append(webParam.HeaderParams, param)
	case "body":
		if webParam.BodyParam != nil {
			utils.Failuref("%s %s, body cannot define multiple", commonFunc.Loc.String(), param.Comment)
		}
		webParam.BodyParam = param
	case "formData":
		webParam.FormParams = append(webParam.FormParams, param)
	default:
		utils.Failuref("%s %s, %s can not support", commonFunc.Loc.String(), param.Comment, param.Source)
	}
}
