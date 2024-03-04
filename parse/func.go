package parse

import (
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

// FuncParse
// 解析方法 -> 注解 -> 生成代码: 当前代码
func (p *Parser) FuncParse(funcDecl *ast.FuncDecl, packageInfo *model.PackageInfo) {
	funcName := funcDecl.Name.String()

	funcInfo := model.NewFuncInfo()
	funcInfo.PackageInfo = packageInfo
	funcInfo.FuncName = funcName
	funcInfo.Proxy = funcName

	astRec := funcDecl.Recv

	if astRec != nil {
		fieldRec := astRec.List[0]
		paramInfo := utils.ToFileInfo(fieldRec)
		funcInfo.Recv = paramInfo
		funcInfo.Params = append(funcInfo.Params, funcInfo.Recv)
	}

	fillEmptyParam(funcDecl.Type, funcInfo)

	isWebApp := false
	isMiddleware := false
	isRouter := false
	hasAnnotate := false
	if funcDecl.Doc != nil {
		for _, comment := range funcDecl.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@proxy" {
				if argsLen >= 2 {
					proxy := annotateArgs[1]
					if proxy != "" && proxy != "_" {
						funcInfo.Proxy = proxy
					}
				}
				funcInfo.ProxyComment = comment.Text
				hasAnnotate = true
			} else if annotateName == "@import" {
				importInfo := &model.ImportInfo{}
				funcInfo.Imports = append(funcInfo.Imports, importInfo)

				if argsLen < 2 {
					utils.Failuref("%s, Path must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				importInfo.Path = annotateArgs[1]
				if argsLen >= 3 {
					importName := annotateArgs[2]
					if importName == "." {
						utils.Failuref("%s, Cannot support DotImport, at %s()", comment.Text, funcInfo.FuncName)
					}
					if importName != "" {
						importInfo.Name = importName
					}
				}
			} else if annotateName == "@injectParam" {
				if argsLen < 2 {
					utils.Failuref("%s, ParamName must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				paramName := annotateArgs[1]
				paramInfo := utils.FindParamInfo(funcInfo, paramName)
				if paramInfo == nil {
					utils.Failuref("%s, ParamName not found, at %s()", comment.Text, funcInfo.FuncName)
				}

				if argsLen >= 3 {
					paramInstance := annotateArgs[2]
					if paramInstance != "" && paramInstance != "_" {
						paramInfo.Instance = paramInstance
					}
				}
				paramInfo.Comment = comment.Text
				paramInfo.Source = "inject"
			} else if annotateName == "@injectRecv" {
				if argsLen < 2 {
					utils.Failuref("%s, RecvName must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				paramName := annotateArgs[1]
				if funcInfo.Recv.Name != paramName {
					utils.Failuref("%s, RecvName not found, at %s()", comment.Text, funcInfo.FuncName)
				}

				paramInfo := funcInfo.Recv
				if argsLen >= 3 {
					paramInstance := annotateArgs[2]
					if paramInstance != "" && paramInstance != "_" {
						paramInfo.Instance = paramInstance
					}
				}
				paramInfo.Comment = comment.Text
				paramInfo.Source = "inject"
			} else if annotateName == "@injectCtx" {
				if argsLen < 2 {
					utils.Failuref("%s, ParamName must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				paramName := annotateArgs[1]
				paramInfo := utils.FindParamInfo(funcInfo, paramName)
				if paramInfo == nil {
					utils.Failuref("%s, ParamName not found, at %s()", comment.Text, funcInfo.FuncName)
				}

				paramInfo.Comment = comment.Text
				paramInfo.Source = "ctx"
			} else if annotateName == "@webAppProvide" {
				isWebApp = true
				hasAnnotate = true
				if isMiddleware {
					utils.Failuref("%s, conflict with %s, at %s()", comment.Text, "@middleware", funcInfo.FuncName)
				} else if isRouter {
					utils.Failuref("%s, conflict with %s, at %s()", comment.Text, "@router", funcInfo.FuncName)
				}
			} else if annotateName == "@middleware" {
				isMiddleware = true
				hasAnnotate = true
				if isWebApp {
					utils.Failuref("%s, conflict with %s, at %s()", comment.Text, "@webAppProvide", funcInfo.FuncName)
				} else if isRouter {
					utils.Failuref("%s, conflict with %s, at %s()", comment.Text, "@router", funcInfo.FuncName)
				}
			} else if annotateName == "@router" {
				isRouter = true
				hasAnnotate = true
				if isWebApp {
					utils.Failuref("%s, conflict with %s, at %s()", comment.Text, "@webAppProvide", funcInfo.FuncName)
				} else if isMiddleware {
					utils.Failuref("%s, conflict with %s, at %s()", comment.Text, "@middleware", funcInfo.FuncName)
				}
			}
		}
	}

	if !hasAnnotate {
		return
	}

	results := funcDecl.Type.Results
	if results != nil {
		for _, resultField := range results.List {
			funcInfo.Results = append(funcInfo.Results, utils.ToFileInfo(resultField))
		}
	}

	for _, param := range funcDecl.Type.Params.List {
		ParamParse(param, funcInfo)
	}

	if isWebApp {
		p.WebAppParse(funcDecl, funcInfo)
	} else if isMiddleware {
		p.MiddlewareParse(funcDecl, funcInfo)
	} else if isRouter {
		p.RouterParse(funcDecl, funcInfo)
	} else {
		addFuncOrMethodInstances(p.Result, funcInfo)
	}
}
func fillEmptyParam(funcType *ast.FuncType, funcInfo *model.FuncInfo) {
	for _, field := range funcType.Params.List {
		funcInfo.Params = append(funcInfo.Params, utils.ToFileInfo(field))
	}
}

func addFuncOrMethodInstances(result *model.ModuleInfo, funcInfo *model.FuncInfo) {
	if funcInfo.Recv == nil {
		result.FuncInstances = append(result.FuncInstances, funcInfo)
	} else {
		result.MethodInstances = append(result.MethodInstances, funcInfo)
	}
}
