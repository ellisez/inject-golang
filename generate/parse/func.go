package parse

import (
	"github.com/ellisez/inject-golang/generate/model"
	"github.com/ellisez/inject-golang/generate/utils"
	"go/ast"
)

// FuncParse
// 解析方法 -> 注解 -> 生成代码: 当前代码
func (p *Parser) FuncParse(funcDecl *ast.FuncDecl, packageInfo *model.PackageInfo) {
	hasAnnotate := false
	funcInfo := &model.FuncInfo{PackageInfo: packageInfo}
	if funcDecl.Doc != nil {
		for _, comment := range funcDecl.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@inject" {
				if argsLen >= 1 {
					funcInfo.Proxy = annotateArgs[1]
				}
				hasAnnotate = true
			}
		}
	}

	if !hasAnnotate {
		return
	}

	astRec := funcDecl.Recv

	if astRec != nil {
		fieldRec := astRec.List[0]
		paramInfo := &model.ParamInfo{}
		paramInfo.Name = fieldRec.Names[0].String()
		paramInfo.Type = utils.TypeToString(fieldRec.Type)
		funcInfo.Recv = paramInfo
	}

	results := funcDecl.Type.Results
	if results != nil {
		resultLen := len(results.List)
		funcInfo.Results = make([]*model.ParamInfo, resultLen)
		for i, resultField := range results.List {
			paramInfo := &model.ParamInfo{}
			if resultField.Names != nil {
				paramInfo.Name = resultField.Names[0].String()
			}
			paramInfo.Type = utils.TypeToString(resultField.Type)
			funcInfo.Results[i] = paramInfo
		}
	}
	funcInfo.FuncName = funcDecl.Name.String()
	for _, param := range funcDecl.Type.Params.List {
		ParamParse(param, funcInfo)
	}

	addFuncOrMethodInstances(p.Result, funcInfo)
}

func addFuncOrMethodInstances(result *model.AnnotateInfo, funcInfo *model.FuncInfo) {
	if funcInfo.Recv == nil {
		addFuncInstances(result, funcInfo)
	} else {
		addMethodInstances(result, funcInfo)
	}
}

func addFuncInstances(result *model.AnnotateInfo, funcInfo *model.FuncInfo) {
	if result.FuncInstances == nil {
		result.FuncInstances = make([]*model.FuncInfo, 0)
	}
	result.FuncInstances = append(result.FuncInstances, funcInfo)
}

func addMethodInstances(result *model.AnnotateInfo, funcInfo *model.FuncInfo) {
	if result.MethodInstances == nil {
		result.MethodInstances = make([]*model.FuncInfo, 0)
	}
	result.MethodInstances = append(result.MethodInstances, funcInfo)
}
