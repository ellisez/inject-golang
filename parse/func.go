package parse

import (
	"fmt"
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
		addParam(funcInfo, funcInfo.Recv)
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
					utils.Failure(fmt.Sprintf("%s, Path must be specified", comment.Text))
				}
				importInfo.Path = annotateArgs[1]
				if argsLen >= 3 {
					importName := annotateArgs[2]
					if importName == "." {
						utils.Failure(fmt.Sprintf("%s, Cannot support DotImport", comment.Text))
					}
					if importName != "" {
						importInfo.Name = importName
					}
				}
			} else if annotateName == "@injectParam" {
				if argsLen < 2 {
					utils.Failure(fmt.Sprintf("%s, ParamName must be specified", comment.Text))
				}
				paramName := annotateArgs[1]
				paramInfo := utils.FindParamInfo(funcInfo, paramName)
				if paramInfo == nil {
					utils.Failure(fmt.Sprintf("%s, ParamName not found", comment.Text))
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
					utils.Failure(fmt.Sprintf("%s, RecvName must be specified", comment.Text))
				}
				paramName := annotateArgs[1]
				if funcInfo.Recv.Name != paramName {
					utils.Failure(fmt.Sprintf("%s, RecvName not found", comment.Text))
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
			} else if annotateName == "@webAppProvide" {
				isWebApp = true
				hasAnnotate = true
				if isMiddleware {
					utils.Failure(fmt.Sprintf("%s, conflict with %s", comment.Text, "@middleware"))
				} else if isRouter {
					utils.Failure(fmt.Sprintf("%s, conflict with %s", comment.Text, "@router"))
				}
			} else if annotateName == "@middleware" {
				isMiddleware = true
				hasAnnotate = true
				if isWebApp {
					utils.Failure(fmt.Sprintf("%s, conflict with %s", comment.Text, "@webAppProvide"))
				} else if isRouter {
					utils.Failure(fmt.Sprintf("%s, conflict with %s", comment.Text, "@router"))
				}
			} else if annotateName == "@router" {
				isRouter = true
				hasAnnotate = true
				if isWebApp {
					utils.Failure(fmt.Sprintf("%s, conflict with %s", comment.Text, "@webAppProvide"))
				} else if isMiddleware {
					utils.Failure(fmt.Sprintf("%s, conflict with %s", comment.Text, "@middleware"))
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
func addParam(funcInfo *model.FuncInfo, paramInfo *model.FieldInfo) {
	if funcInfo.Params != nil {
		funcInfo.Params = make([]*model.FieldInfo, 0)
	}
	funcInfo.Params = append(funcInfo.Params, paramInfo)
}
func fillEmptyParam(funcType *ast.FuncType, funcInfo *model.FuncInfo) {
	for _, field := range funcType.Params.List {
		addParam(funcInfo, utils.ToFileInfo(field))
	}
}

func addFuncOrMethodInstances(result *model.ModuleInfo, funcInfo *model.FuncInfo) {
	if funcInfo.Recv == nil {
		result.FuncInstances = append(result.FuncInstances, funcInfo)
	} else {
		result.MethodInstances = append(result.MethodInstances, funcInfo)
	}
}
