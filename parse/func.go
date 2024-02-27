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
	funcInfo := &model.FuncInfo{
		PackageInfo: packageInfo,
		FuncName:    funcName,
		Proxy:       funcName,
	}
	fillEmptyParam(funcDecl.Type, funcInfo)

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
				if funcInfo.Imports == nil {
					funcInfo.Imports = []*model.ImportInfo{
						importInfo,
					}
				} else {
					funcInfo.Imports = append(funcInfo.Imports, importInfo)
				}

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
				paramInfo.IsInject = true
			}
		}
	}

	if !hasAnnotate {
		return
	}

	astRec := funcDecl.Recv

	if astRec != nil {
		fieldRec := astRec.List[0]
		paramInfo := &model.FieldInfo{}
		paramInfo.Name = fieldRec.Names[0].String()
		paramInfo.Type = fieldRec.Type
		funcInfo.Recv = paramInfo
	}

	results := funcDecl.Type.Results
	if results != nil {
		resultLen := len(results.List)
		funcInfo.Results = make([]*model.FieldInfo, resultLen)
		for i, resultField := range results.List {
			paramInfo := &model.FieldInfo{}
			if resultField.Names != nil {
				paramInfo.Name = resultField.Names[0].String()
			}
			paramInfo.Type = resultField.Type
			funcInfo.Results[i] = paramInfo
		}
	}

	for _, param := range funcDecl.Type.Params.List {
		ParamParse(param, funcInfo)
	}

	addFuncOrMethodInstances(p.Result, funcInfo)
}

func fillEmptyParam(funcType *ast.FuncType, funcInfo *model.FuncInfo) {
	fl := 0
	if funcType.Params != nil {
		fl = len(funcType.Params.List)
	}
	funcInfo.Params = make([]*model.FieldInfo, fl)
	for i, field := range funcType.Params.List {
		funcInfo.Params[i] = utils.ToFileInfo(field)
	}
}

func addFuncOrMethodInstances(result *model.ModuleInfo, funcInfo *model.FuncInfo) {
	if funcInfo.Recv == nil {
		addFuncInstances(result, funcInfo)
	} else {
		addMethodInstances(result, funcInfo)
	}
}

func addFuncInstances(result *model.ModuleInfo, funcInfo *model.FuncInfo) {
	if result.FuncInstances == nil {
		result.FuncInstances = make([]*model.FuncInfo, 0)
	}
	result.FuncInstances = append(result.FuncInstances, funcInfo)
}

func addMethodInstances(result *model.ModuleInfo, funcInfo *model.FuncInfo) {
	if result.MethodInstances == nil {
		result.MethodInstances = make([]*model.FuncInfo, 0)
	}
	result.MethodInstances = append(result.MethodInstances, funcInfo)
}
