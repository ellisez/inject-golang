package parse

import (
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"strconv"
	"strings"
)

func (p *Parser) WebAppParse(funcDecl *ast.FuncDecl, funcInfo *model.FuncInfo) {
	webInfo := model.NewWebInfoFromFunc(funcInfo)
	if funcDecl.Doc != nil {
		for _, comment := range funcDecl.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@webAppProvide" {
				if argsLen >= 2 {
					webInfo.WebApp = annotateArgs[1]
				}
				webInfo.WebAppComment = comment.Text
			} else if annotateName == "@static" {
				staticResource := model.NewStaticResource()
				webInfo.Statics = append(webInfo.Statics, staticResource)

				if argsLen < 2 {
					utils.Failuref("%s, Path must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				staticResource.Path = annotateArgs[1]
				if argsLen < 3 {
					utils.Failuref("%s, Dirname must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				staticResource.Path = annotateArgs[2]

				if argsLen >= 4 {
					features := annotateArgs[3]
					if !strings.HasPrefix(features, "[") || !strings.HasSuffix(features, "]") {
						utils.Failuref("%s, Features must be wrapped in [], at %s()", comment.Text, funcInfo.FuncName)
					}
					features = features[1 : len(features)-1]
					splitStr := strings.Split(features, ",")
					for _, feature := range splitStr {
						staticResource.Features = append(staticResource.Features, strings.TrimSpace(feature))
					}
				}
				if argsLen >= 5 {
					staticResource.Index = annotateArgs[4]
				}
				if argsLen >= 6 {
					maxAge, err := strconv.Atoi(annotateArgs[5])
					if err != nil {
						utils.Failuref("%s, MaxAge must be a number, at %s()", comment.Text, funcInfo.FuncName)
					}
					staticResource.MaxAge = maxAge
				}
				staticResource.StaticComment = comment.Text
			} else if annotateName == "@proxy" {
				if argsLen >= 2 {
					proxy := annotateArgs[1]
					if proxy != "" && proxy != "_" {
						funcInfo.Proxy = proxy
					}
				}
				funcInfo.ProxyComment = comment.Text
			}
		}
	}

	createdWebApp := p.Result.GetWebApp(webInfo.WebApp)
	if createdWebApp != nil {
		createdWebApp.WebAppComment = webInfo.WebAppComment
		createdWebApp.FuncInfo = webInfo.FuncInfo
		createdWebApp.Statics = append(createdWebApp.Statics, webInfo.Statics...)
	} else {
		p.Result.WebAppInstances = append(p.Result.WebAppInstances, webInfo)
	}
}

func addRouterParam(routerParam *model.RouterParam, paramInfo *model.FieldInfo, funcInfo *model.FuncInfo) {
	switch paramInfo.Source {
	case "query":
		routerParam.QueryParams = append(routerParam.QueryParams, paramInfo)
		break
	case "path":
		routerParam.PathParams = append(routerParam.PathParams, paramInfo)
		break
	case "header":
		routerParam.HeaderParams = append(routerParam.HeaderParams, paramInfo)
		break
	case "body":
		if routerParam.BodyParam != nil {
			utils.Failuref("%s, body cannot define multiple, at %s()", paramInfo.Comment, funcInfo.FuncName)
		}
		routerParam.BodyParam = paramInfo
		break
	case "formData":
		routerParam.FormParams = append(routerParam.FormParams, paramInfo)
		break
	default:
		utils.Failuref("%s, %s can not support, at %s()", paramInfo.Comment, paramInfo.Source, funcInfo.FuncName)
	}
}

func (p *Parser) MiddlewareParse(funcDecl *ast.FuncDecl, funcInfo *model.FuncInfo) {
	middlewareInfo := model.NewMiddlewareInfoFromFuncInfo(funcInfo)

	if funcDecl.Doc != nil {
		for _, comment := range funcDecl.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@middleware" {
				if argsLen < 2 {
					utils.Failuref("%s, Path must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				middlewareInfo.Path = annotateArgs[1]

				middlewareInfo.MiddlewareComment = comment.Text
			} else if annotateName == "@param" {
				if argsLen < 2 {
					utils.Failuref("%s, ParamName must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				paramName := annotateArgs[1]
				paramInfo := utils.FindParamInfo(funcInfo, paramName)
				if paramInfo == nil {
					utils.Failuref("%s, ParamName not found, at %s()", comment.Text, funcInfo.FuncName)
				}

				if argsLen < 3 {
					utils.Failuref("%s, ParamName must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				paramSource := annotateArgs[2]
				paramInfo.Comment = comment.Text
				paramInfo.Source = paramSource
				addRouterParam(middlewareInfo.RouterParam, paramInfo, funcInfo)
			} else if annotateName == "@injectWebCtx" {
				if argsLen < 2 {
					utils.Failuref("%s, ParamName must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				paramName := annotateArgs[1]
				paramInfo := utils.FindParamInfo(funcInfo, paramName)
				if paramInfo == nil {
					utils.Failuref("%s, ParamName not found, at %s()", comment.Text, funcInfo.FuncName)
				}

				paramInfo.Comment = comment.Text
				paramInfo.Source = "webCtx"
			} else if annotateName == "@webApp" {
				if argsLen >= 2 {
					middlewareInfo.WebApp = annotateArgs[1]
				}
			}
		}
	}

	for _, param := range middlewareInfo.Params {
		if param.Source == "" {
			utils.Failuref("%s, The ParamName \"%s\" does not allow non injection, at %s()", middlewareInfo.MiddlewareComment, param.Name, funcInfo.FuncName)
		}
	}

	createdWebApp := p.Result.GetWebApp(middlewareInfo.WebApp)
	if createdWebApp != nil {
		if createdWebApp.Middlewares == nil {
			createdWebApp.Middlewares = []*model.MiddlewareInfo{
				middlewareInfo,
			}
		} else {
			createdWebApp.Middlewares = append(createdWebApp.Middlewares, middlewareInfo)
		}
	} else {
		webInfo := model.NewWebInfo()
		webInfo.Middlewares = []*model.MiddlewareInfo{
			middlewareInfo,
		}
		if p.Result.WebAppInstances == nil {
			p.Result.WebAppInstances = []*model.WebInfo{
				webInfo,
			}
		} else {
			p.Result.WebAppInstances = append(p.Result.WebAppInstances, webInfo)
		}
	}
}

func (p *Parser) RouterParse(funcDecl *ast.FuncDecl, funcInfo *model.FuncInfo) {
	routerInfo := model.NewRouterInfoFromFuncInfo(funcInfo)

	if funcDecl.Doc != nil {
		for _, comment := range funcDecl.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@router" {
				if argsLen < 2 {
					utils.Failuref("%s, Path must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				routerInfo.Path = annotateArgs[1]

				if argsLen < 3 {
					utils.Failuref("%s, Methods must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				methods := annotateArgs[2]
				if !strings.HasPrefix(methods, "[") || !strings.HasSuffix(methods, "]") {
					utils.Failuref("%s, Methods must be wrapped in [], at %s()", comment.Text, funcInfo.FuncName)
				}
				methods = methods[1 : len(methods)-1]
				for _, method := range strings.Split(methods, ",") {
					routerInfo.Methods = append(routerInfo.Methods, utils.FirstToUpper(method))
				}

				routerInfo.RouterComment = comment.Text
			} else if annotateName == "@param" {
				if argsLen < 2 {
					utils.Failuref("%s, ParamName must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				paramName := annotateArgs[1]
				paramInfo := utils.FindParamInfo(funcInfo, paramName)
				if paramInfo == nil {
					utils.Failuref("%s, ParamName not found, at %s()", comment.Text, funcInfo.FuncName)
				}

				if argsLen < 3 {
					utils.Failuref("%s, ParamName must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				paramSource := annotateArgs[2]
				paramInfo.Comment = comment.Text
				if paramInfo.Source != "" {
					utils.Failuref("%s, conflict with %s, at %s()", comment.Text, paramInfo.Source, funcInfo.FuncName)
				}
				paramInfo.Source = paramSource
				addRouterParam(routerInfo.RouterParam, paramInfo, funcInfo)
			} else if annotateName == "@webApp" {
				if argsLen >= 2 {
					routerInfo.WebApp = annotateArgs[1]
				}
			}
		}
	}

	for _, param := range routerInfo.Params {
		if param.Source == "" {
			utils.Failuref("%s, The ParamName \"%s\" does not allow non injection, at %s()", routerInfo.RouterComment, param.Name, funcInfo.FuncName)
		}
	}

	createdWebApp := p.Result.GetWebApp(routerInfo.WebApp)
	if createdWebApp != nil {
		if createdWebApp.Routers == nil {
			createdWebApp.Routers = []*model.RouterInfo{
				routerInfo,
			}
		} else {
			createdWebApp.Routers = append(createdWebApp.Routers, routerInfo)
		}
	} else {
		webInfo := model.NewWebInfo()
		webInfo.Routers = []*model.RouterInfo{
			routerInfo,
		}
		if p.Result.WebAppInstances == nil {
			p.Result.WebAppInstances = []*model.WebInfo{
				webInfo,
			}
		} else {
			p.Result.WebAppInstances = append(p.Result.WebAppInstances, webInfo)
		}
	}
}
