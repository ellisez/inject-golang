package parse

import (
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

// FieldParse
// 解析属性 -> 注解 -> 生成代码: 当前代码
func FieldParse(field *ast.Field, structInfo *model.StructInfo) {
	fieldName := utils.FieldName(field)

	fieldInfo := utils.FindFieldInfo(structInfo, fieldName)

	if field.Doc != nil {
		for _, comment := range field.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}

			annotateName := annotateArgs[0]
			if annotateName == "@inject" {
				if argsLen >= 2 {
					fieldInstance := annotateArgs[1]
					if fieldInstance != "" && fieldInstance != "_" {
						fieldInfo.Instance = fieldInstance
					}
				}
				if argsLen >= 3 {
					fieldGetter := annotateArgs[2]
					if fieldGetter != "" && fieldGetter != "_" {
						fieldInfo.Getter = fieldGetter
					}
				}

				if argsLen >= 4 {
					fieldSetter := annotateArgs[3]
					if fieldSetter != "" && fieldSetter != "_" {
						fieldInfo.Setter = fieldSetter
					}
				}
				fieldInfo.Comment = comment.Text
				fieldInfo.Source = "inject"
			}
		}
	}

}
