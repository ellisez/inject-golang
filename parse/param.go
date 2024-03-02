package parse

import (
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

// ParamParse
// 解析参数 -> 注解 -> 生成代码: 当前代码
func ParamParse(param *ast.Field, funcInfo *model.FuncInfo) {
	var paramName string
	if param.Names != nil {
		paramName = param.Names[0].String()
	}

	paramInfo := utils.FindParamInfo(funcInfo, paramName)

	paramInfo.Type = param.Type

	if param.Doc != nil {
		for _, comment := range param.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@inject" {
				if argsLen >= 2 {
					paramInstance := annotateArgs[1]
					if paramInstance != "" && paramInstance != "_" {
						paramInfo.Instance = paramInstance
					}
				}
				paramInfo.Comment = comment.Text
				paramInfo.Source = "inject"
			}
		}
	}
}
