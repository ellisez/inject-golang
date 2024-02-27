package parse

import (
	"fmt"
	"github.com/ellisez/inject-golang/generate/model"
	"github.com/ellisez/inject-golang/generate/utils"
	"go/ast"
)

// StructParse
// 解析结构体 -> 注解 -> 生成代码: 全局与当前
func (p *Parser) StructParse(structDecl *ast.GenDecl, packageInfo *model.PackageInfo) {
	typeSpec := structDecl.Specs[0].(*ast.TypeSpec)
	structType := typeSpec.Type.(*ast.StructType)
	structName := typeSpec.Name.String()
	structInfo := &model.StructInfo{
		PackageInfo: packageInfo,
		Name:        structName,
		Instance:    structName,
	}

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
				if structInfo.Imports == nil {
					structInfo.Imports = []*model.ImportInfo{
						importInfo,
					}
				} else {
					structInfo.Imports = append(structInfo.Imports, importInfo)
				}

				if argsLen < 2 {
					panic(fmt.Errorf("%s, Path must be specified", comment.Text))
				}
				importInfo.Path = annotateArgs[1]

				if argsLen >= 3 {
					importName := annotateArgs[2]
					if importName == "." {
						panic(fmt.Errorf("%s, Cannot support DotImport", comment.Text))
					}
					if importName != "" {
						importInfo.Name = importName
					}
				}
			} else if annotateName == "@injectField" {
				if argsLen < 2 {
					panic(fmt.Errorf("%s, FieldName must be specified", comment.Text))
				}
				fieldName := annotateArgs[1]
				fieldInfo := utils.FindFieldInfo(structInfo, fieldName)
				if fieldInfo == nil {
					panic(fmt.Errorf("%s, FieldName name not found", comment.Text))
				}
				fieldInfo.Comment = comment.Text

				if argsLen >= 3 {
					fieldInstance := annotateArgs[2]
					if fieldInstance != "" && fieldInstance != "_" {
						fieldInfo.Instance = fieldInstance
					}
				}
				fieldInfo.IsInject = true
			} else if annotateName == "@preConstruct" {
				if argsLen < 2 {
					panic(fmt.Errorf("%s, FuncName must be specified", comment.Text))
				}
				structInfo.PreConstructComment = comment.Text
				structInfo.PreConstruct = annotateArgs[1]
			} else if annotateName == "@postConstruct" {
				if argsLen < 2 {
					panic(fmt.Errorf("%s, FuncName must be specified", comment.Text))
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
	fl := 0
	if structType.Fields != nil {
		fl = len(structType.Fields.List)
	}
	structInfo.Fields = make([]*model.FieldInfo, fl)
	for i, field := range structType.Fields.List {
		structInfo.Fields[i] = utils.ToFileInfo(field)
	}
}

func addStructInstances(result *model.ModuleInfo, structInfo *model.StructInfo) {
	if structInfo.Mode == "multiple" {
		addMultipleInstances(result, structInfo)
	} else {
		addSingletonInstances(result, structInfo)
	}
}

func addSingletonInstances(result *model.ModuleInfo, structInfo *model.StructInfo) {
	if result.SingletonInstances == nil {
		result.SingletonInstances = make([]*model.StructInfo, 0)
	}
	result.SingletonInstances = append(result.SingletonInstances, structInfo)
}

func addMultipleInstances(result *model.ModuleInfo, structInfo *model.StructInfo) {
	if result.MultipleInstances == nil {
		result.MultipleInstances = make([]*model.StructInfo, 0)
	}
	result.MultipleInstances = append(result.MultipleInstances, structInfo)
}
