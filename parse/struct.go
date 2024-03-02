package parse

import (
	"fmt"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

// StructParse
// 解析结构体 -> 注解 -> 生成代码: 全局与当前
func (p *Parser) StructParse(structDecl *ast.GenDecl, packageInfo *model.PackageInfo) {
	typeSpec := structDecl.Specs[0].(*ast.TypeSpec)
	structType := typeSpec.Type.(*ast.StructType)
	structName := typeSpec.Name.String()

	structInfo := model.NewStructInfo()
	structInfo.PackageInfo = packageInfo
	structInfo.Name = structName
	structInfo.Instance = structName

	fillEmptyField(structType, structInfo)
	hasAnnotate := false

	if structDecl.Doc != nil {
		for _, comment := range structDecl.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@provide" {
				if argsLen >= 2 {
					structInstance := annotateArgs[1]
					if structInstance != "" && structInstance != "_" {
						structInfo.Instance = structInstance
					}
				}
				if argsLen >= 3 {
					structInfo.Mode = annotateArgs[2]
				}
				structInfo.ProvideComment = comment.Text
				hasAnnotate = true
			} else if annotateName == "@import" {
				importInfo := &model.ImportInfo{}
				structInfo.Imports = append(structInfo.Imports, importInfo)

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
			} else if annotateName == "@injectField" {
				if argsLen < 2 {
					utils.Failure(fmt.Sprintf("%s, FieldName must be specified", comment.Text))
				}
				fieldName := annotateArgs[1]
				fieldInfo := utils.FindFieldInfo(structInfo, fieldName)
				if fieldInfo == nil {
					utils.Failure(fmt.Sprintf("%s, FieldName name not found", comment.Text))
				}
				fieldInfo.Comment = comment.Text

				if argsLen >= 3 {
					fieldInstance := annotateArgs[2]
					if fieldInstance != "" && fieldInstance != "_" {
						fieldInfo.Instance = fieldInstance
					}
				}
				fieldInfo.Source = "inject"
			} else if annotateName == "@preConstruct" {
				if argsLen < 2 {
					utils.Failure(fmt.Sprintf("%s, FuncName must be specified", comment.Text))
				}
				structInfo.PreConstructComment = comment.Text
				structInfo.PreConstruct = annotateArgs[1]
			} else if annotateName == "@postConstruct" {
				if argsLen < 2 {
					utils.Failure(fmt.Sprintf("%s, FuncName must be specified", comment.Text))
				}
				structInfo.PostConstructComment = comment.Text
				structInfo.PostConstruct = annotateArgs[1]
			}
		}
	}

	if !hasAnnotate {
		return
	}

	for _, field := range structType.Fields.List {
		FieldParse(field, structInfo)
	}

	addStructInstances(p.Result, structInfo)
}

func fillEmptyField(structType *ast.StructType, structInfo *model.StructInfo) {
	for _, field := range structType.Fields.List {
		structInfo.Fields = append(structInfo.Fields, utils.ToFileInfo(field))
	}
}

func addStructInstances(result *model.ModuleInfo, structInfo *model.StructInfo) {
	if structInfo.Mode == "multiple" {
		result.SingletonInstances = append(result.SingletonInstances, structInfo)
	} else {
		result.MultipleInstances = append(result.MultipleInstances, structInfo)
	}
}
