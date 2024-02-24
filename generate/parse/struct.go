package parse

import (
	"github.com/ellisez/inject-golang/generate/model"
	"go/ast"
)

// StructParse
// 解析结构体 -> 注解 -> 生成代码: 全局与当前
func (p *Parser) StructParse(structDecl *ast.GenDecl, packageInfo *model.PackageInfo) {
	hasAnnotate := false
	structInfo := &model.StructInfo{PackageInfo: packageInfo}
	if structDecl.Doc != nil {
		for _, comment := range structDecl.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@provide" {
				if argsLen >= 1 {
					structInfo.Instance = annotateArgs[1]
				}
				if argsLen >= 2 {
					structInfo.Mode = annotateArgs[2]
				}
				hasAnnotate = true
			} else if annotateName == "@constructor" {
				if argsLen >= 1 {
					structInfo.Constructor = annotateArgs[1]
				}
				hasAnnotate = true
			}
		}
	}

	if !hasAnnotate {
		return
	}

	typeSpec := structDecl.Specs[0].(*ast.TypeSpec)
	structInfo.Name = typeSpec.Name.String()

	structType := typeSpec.Type.(*ast.StructType)
	for _, field := range structType.Fields.List {
		FieldParse(field, structInfo)
	}

	addStructInstances(p.Result, structInfo)
}

func addStructInstances(result *model.AnnotateInfo, structInfo *model.StructInfo) {
	switch structInfo.Mode {
	case "singleton":
		addSingletonInstances(result, structInfo)
		break
	case "multiple":
		addMultipleInstances(result, structInfo)
		break
	}
}

func addSingletonInstances(result *model.AnnotateInfo, structInfo *model.StructInfo) {
	if result.SingletonInstances == nil {
		result.SingletonInstances = make([]*model.StructInfo, 0)
	}
	result.SingletonInstances = append(result.SingletonInstances, structInfo)
}

func addMultipleInstances(result *model.AnnotateInfo, structInfo *model.StructInfo) {
	if result.MultipleInstances == nil {
		result.MultipleInstances = make([]*model.StructInfo, 0)
	}
	result.MultipleInstances = append(result.MultipleInstances, structInfo)
}
