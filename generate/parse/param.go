package parse

import (
	"github.com/ellisez/inject-golang/generate/model"
	"github.com/ellisez/inject-golang/generate/utils"
	"go/ast"
)

// ParamParse
// 解析参数 -> 注解 -> 生成代码: 当前代码
func ParamParse(param *ast.Field, funcInfo *model.FuncInfo) {
	var paramName string
	if param.Names != nil {
		paramName = param.Names[0].String()
	}

	paramInfo := findParam(funcInfo, paramName)
	if paramInfo == nil {
		paramInfo = &model.ParamInfo{
			Name: paramName,
		}
		if funcInfo.Params == nil {
			funcInfo.Params = make([]*model.ParamInfo, 0)
		}
		funcInfo.Params = append(funcInfo.Params, paramInfo)
	}

	paramInfo.Type = utils.TypeToString(param.Type)

	if param.Doc != nil {
		for _, comment := range param.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@inject" {
				if argsLen >= 1 {
					paramInfo.Instance = annotateArgs[1]
					paramInfo.IsInject = true
				}
			}
		}
	}

	if paramInfo.IsInject {
		addInjectParam(funcInfo, paramInfo)
	} else {
		addNormalParam(funcInfo, paramInfo)
	}
}

func findParam(funcInfo *model.FuncInfo, paramName string) *model.ParamInfo {
	var paramInfo *model.ParamInfo
	for _, funcParam := range funcInfo.Params {
		if funcParam.Name == paramName {
			paramInfo = funcParam
			break
		}
	}
	return paramInfo
}
func addInjectParam(funcInfo *model.FuncInfo, paramInfo *model.ParamInfo) {
	if funcInfo.InjectParams == nil {
		funcInfo.InjectParams = make([]*model.ParamInfo, 0)
	}
	funcInfo.InjectParams = append(funcInfo.InjectParams, paramInfo)
}

func addNormalParam(funcInfo *model.FuncInfo, paramInfo *model.ParamInfo) {
	if funcInfo.NormalParams == nil {
		funcInfo.NormalParams = make([]*model.ParamInfo, 0)
	}
	funcInfo.NormalParams = append(funcInfo.NormalParams, paramInfo)
}
