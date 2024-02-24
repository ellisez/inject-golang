package parse

import (
	"github.com/ellisez/inject-golang/generate/model"
	"github.com/ellisez/inject-golang/generate/utils"
	"go/ast"
)

// FieldParse
// 解析属性 -> 注解 -> 生成代码: 当前代码
func FieldParse(field *ast.Field, structInfo *model.StructInfo) {
	var fieldName string
	if field.Names != nil {
		fieldName = field.Names[0].String()
	}

	fieldInfo := findField(structInfo, fieldName)
	if fieldInfo == nil {
		fieldInfo = &model.FieldInfo{
			Name: fieldName,
		}
		if structInfo.Fields == nil {
			structInfo.Fields = make([]*model.FieldInfo, 0)
		}
		structInfo.Fields = append(structInfo.Fields, fieldInfo)
	}
	fieldInfo.Type = utils.TypeToString(field.Type)

	hasAnnotate := false
	if field.Doc != nil {
		for _, comment := range field.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}

			annotateName := annotateArgs[0]
			if annotateName == "@inject" {
				if argsLen >= 1 {
					fieldInfo.Instance = annotateArgs[1]
				}
				hasAnnotate = true
			}
		}
	}

	fieldInfo.IsInject = false
	if hasAnnotate || structInfo.Mode == "singleton" {
		fieldInfo.IsInject = true
	}

	if fieldInfo.IsInject {
		addInjectField(structInfo, fieldInfo)
	} else {
		addNormalField(structInfo, fieldInfo)
	}
}
func findField(structInfo *model.StructInfo, fieldName string) *model.FieldInfo {
	var fieldInfo *model.FieldInfo
	for _, structField := range structInfo.Fields {
		if structField.Name == fieldName {
			fieldInfo = structField
			break
		}
	}
	return fieldInfo
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
