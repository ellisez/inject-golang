package parse

import (
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"strconv"
	"strings"
)

func (p *Parser) WebParse(_ *ast.FuncDecl, commonFunc *model.CommonFunc, comments []*model.Comment) {
	webFormatCheck(commonFunc)
	webApp := model.NewWebInstance()
	webApp.CommonFunc = commonFunc

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
						utils.Failuref(`%s %s, Instance "%s" must be capitalized with the first letter`, commonFunc.Loc.String(), webApp.Comment, instance)
					}
					webApp.Instance = instance
				}
			}
			webApp.Comment = comment.Comment
		case "@static":
			if argsLen < 2 {
				utils.Failuref("%s %s, Path must be specified", commonFunc.Loc.String(), comment.Comment)
			}
			path := args[1]
			if argsLen < 3 {
				utils.Failuref("%s %s, Dirname must be specified", commonFunc.Loc.String(), comment.Comment)
			}
			dirname := args[2]

			resource := &model.WebResource{
				Path:    path,
				Dirname: dirname,
			}
			webApp.Resources = append(webApp.Resources, resource)

			if argsLen >= 4 {
				features := args[3]
				if !strings.HasPrefix(features, "[") || !strings.HasSuffix(features, "]") {
					utils.Failuref("%s %s, Features must be wrapped in []", commonFunc.Loc.String(), comment.Comment)
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
					utils.Failuref("%s %s, MaxAge must be a number", commonFunc.Loc.String(), comment.Comment)
				}
				resource.MaxAge = maxAge
			}
			resource.Comment = comment.Comment
		}
	}

	instance := p.Ctx.SingletonOf(webApp.Instance)
	if instance != nil {
		webInstance, ok := instance.(*model.WebInstance)
		if !ok {
			utils.Failuref(`%s %s, Conflict with "%s"`, commonFunc.Loc.String(), webApp.Comment, instance.GetComment())
		}
		webInstance.Comment = webApp.Comment
		webInstance.Imports = append(webInstance.Imports, webApp.Imports...)
		webInstance.Func = webApp.Func
		webInstance.Resources = append(webInstance.Resources, webApp.Resources...)
	} else {
		p.Ctx.SingletonInstances = append(p.Ctx.SingletonInstances, webApp)
	}
	p.Ctx.HasWebInstance = true
}

func webFormatCheck(commonFunc *model.CommonFunc) {
	if len(commonFunc.Results) == 3 {
		utils.Failuref(`%s %s, Illegal webProvide function, the number of returns is not 3`, commonFunc.Loc.String(), commonFunc.Comment)
	}
	host, ok := commonFunc.Results[0].Type.(*ast.Ident)
	if !ok {
		utils.Failuref(`%s %s, Illegal webProvide function, 1st returns is not "string"`, commonFunc.Loc.String(), commonFunc.Comment)
	}
	if host.String() != "string" {
		utils.Failuref(`%s %s, Illegal webProvide function, 1st returns is not "string"`, commonFunc.Loc.String(), commonFunc.Comment)
	}

	port, ok := commonFunc.Results[1].Type.(*ast.Ident)
	if !ok {
		utils.Failuref(`%s %s, Illegal webProvide function, 2st returns is not "uint"`, commonFunc.Loc.String(), commonFunc.Comment)
	}
	if port.String() != "uint" {
		utils.Failuref(`%s %s, Illegal webProvide function, 2st returns is not "uint"`, commonFunc.Loc.String(), commonFunc.Comment)
	}

	err, ok := commonFunc.Results[2].Type.(*ast.Ident)
	if !ok {
		utils.Failuref(`%s %s, Illegal webProvide function, 3st returns is not "error"`, commonFunc.Loc.String(), commonFunc.Comment)
	}
	if err.String() != "error" {
		utils.Failuref(`%s %s, Illegal webProvide function, 3st returns is not "error"`, commonFunc.Loc.String(), commonFunc.Comment)
	}

}

func webParamParse(webParam *model.WebParam, commonFunc *model.CommonFunc, comment *model.Comment) {
	args := comment.Args
	argsLen := len(args)
	if argsLen < 2 {
		utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Comment)
	}
	paramName := args[1]
	paramInfo := utils.FindParam(commonFunc.Func, paramName)
	if paramInfo == nil {
		utils.Failuref("%s %s, ParamName not found", commonFunc.Loc.String(), comment.Comment)
	}

	if argsLen < 3 {
		utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Comment)
	}
	paramSource := args[2]
	paramInfo.Comment = comment.Comment
	paramInfo.Source = paramSource
	addWebParam(webParam, paramInfo, commonFunc)
}

func injectWebCtxParse(commonFunc *model.CommonFunc, comment *model.Comment) {
	args := comment.Args
	argsLen := len(args)
	if argsLen < 2 {
		utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Comment)
	}
	paramName := args[1]
	paramInfo := utils.FindParam(commonFunc.Func, paramName)
	if paramInfo == nil {
		utils.Failuref("%s %s, ParamName not found", commonFunc.Loc.String(), comment.Comment)
	}

	paramInfo.Comment = comment.Comment
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
