package parse

import (
	"fmt"
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
					utils.Failure(fmt.Sprintf("%s, Path must be specified", comment.Text))
				}
				staticResource.Path = annotateArgs[1]
				if argsLen < 3 {
					utils.Failure(fmt.Sprintf("%s, Dirname must be specified", comment.Text))
				}
				staticResource.Path = annotateArgs[2]

				if argsLen >= 4 {
					splitStr := strings.Split(annotateArgs[3], ",")
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
						utils.Failure(fmt.Sprintf("%s, MaxAge must be a number", comment.Text))
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

func addRouterParam(routerParam *model.RouterParam, paramInfo *model.FieldInfo) {
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
			utils.Failure(fmt.Sprintf("%s, body cannot define multiple", paramInfo.Comment))
		}
		routerParam.BodyParam = paramInfo
		break
	case "formData":
		routerParam.FormParams = append(routerParam.FormParams, paramInfo)
		break
	default:
		utils.Failure(fmt.Sprintf("%s, %s can not support", paramInfo.Comment, paramInfo.Source))
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
					utils.Failure(fmt.Sprintf("%s, Path must be specified", comment.Text))
				}
				middlewareInfo.Path = annotateArgs[1]

				middlewareInfo.MiddlewareComment = comment.Text
			} else if annotateName == "@param" {
				if argsLen < 2 {
					utils.Failure(fmt.Sprintf("%s, ParamName must be specified", comment.Text))
				}
				paramName := annotateArgs[1]
				paramInfo := utils.FindParamInfo(funcInfo, paramName)
				if paramInfo == nil {
					utils.Failure(fmt.Sprintf("%s, ParamName not found", comment.Text))
				}

				if argsLen < 3 {
					utils.Failure(fmt.Sprintf("%s, ParamName must be specified", comment.Text))
				}
				paramSource := annotateArgs[2]
				paramInfo.Comment = comment.Text
				paramInfo.Source = paramSource
				addRouterParam(middlewareInfo.RouterParam, paramInfo)
			} else if annotateName == "@webApp" {
				if argsLen >= 2 {
					middlewareInfo.WebApp = annotateArgs[1]
				}
			}
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
	routerInfo := &model.RouterInfo{FuncInfo: funcInfo, WebApp: "WebApp"}

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
					utils.Failure(fmt.Sprintf("%s, Path must be specified", comment.Text))
				}
				routerInfo.Path = annotateArgs[1]

				routerInfo.RouterComment = comment.Text
			} else if annotateName == "@param" {
				if argsLen < 2 {
					utils.Failure(fmt.Sprintf("%s, ParamName must be specified", comment.Text))
				}
				paramName := annotateArgs[1]
				paramInfo := utils.FindParamInfo(funcInfo, paramName)
				if paramInfo == nil {
					utils.Failure(fmt.Sprintf("%s, ParamName not found", comment.Text))
				}

				if argsLen < 3 {
					utils.Failure(fmt.Sprintf("%s, ParamName must be specified", comment.Text))
				}
				paramSource := annotateArgs[2]
				paramInfo.Comment = comment.Text
				if paramInfo.Source != "" {
					utils.Failure(fmt.Sprintf("%s, conflict with %s", comment.Text, paramInfo.Source))
				}
				paramInfo.Source = paramSource
				addRouterParam(routerInfo.RouterParam, paramInfo)
			} else if annotateName == "@webApp" {
				if argsLen >= 2 {
					routerInfo.WebApp = annotateArgs[1]
				}
			}
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
