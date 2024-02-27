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
				fieldInfo.Comment = comment.Text
				fieldInfo.IsInject = true
			}
		}
	}

	if fieldInfo.IsInject {
		addInjectField(structInfo, fieldInfo)
	} else {
		addNormalField(structInfo, fieldInfo)
	}
}
func addInjectField(structInfo *model.StructInfo, fieldInfo *model.FieldInfo) {
	if structInfo.NormalFields == nil {
		structInfo.InjectFields = make([]*model.FieldInfo, 0)
	}
	structInfo.InjectFields = append(structInfo.InjectFields, fieldInfo)
}

func addNormalField(structAnnotate *model.StructInfo, field *model.FieldInfo) {
	if structAnnotate.NormalFields == nil {
		structAnnotate.NormalFields = make([]*model.FieldInfo, 0)
	}
	structAnnotate.NormalFields = append(structAnnotate.NormalFields, field)
}
