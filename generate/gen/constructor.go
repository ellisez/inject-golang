package gen

import (
	"generate/global"
	"generate/model"
	"generate/utils"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
)

// __gen_constructor.go
func genConstructFile(annotateInfo *model.AnnotateInfo, dir string) error {
	filename := filepath.Join(dir, "__gen_container.go")
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	astFile := &ast.File{
		Name:  astIdent(global.GenPackage),
		Scope: ast.NewScope(nil),
	}

	genInjectImportsAst(annotateInfo, astFile)

	genInjectAst(annotateInfo, astFile)

	err = format.Node(file, token.NewFileSet(), astFile)
	if err != nil {
		return err
	}
	return nil
}

func genInjectImportsAst(annotateInfo *model.AnnotateInfo, astFile *ast.File) {

	for _, instance := range annotateInfo.MultipleInstances {
		astImport(astFile, "", instance.Dirname)
	}
}

// # gen segment: Multiple instance #
func genInjectAst(annotateInfo *model.AnnotateInfo, astFile *ast.File) {

	for _, instance := range annotateInfo.MultipleInstances {
		recvVar := utils.FirstToLower(global.StructName)
		param := make([]*ast.Field, 0)
		for _, field := range instance.NormalFields {
			// [code] {{FieldInstance}} {{FieldType}},
			fieldInstance := field.Instance
			if fieldInstance == "" {
				fieldInstance = field.Name
				if field.Name == "" {
					fieldInstance = utils.ShortType(field.Type)
				}
			}
			param = append(param, astField(fieldInstance, utils.TypeToAst(field.Type)))
		}

		provideInstance := instance.Instance
		if provideInstance == "" {
			provideInstance = instance.Name
		}

		instanceVar := utils.FirstToLower(provideInstance)
		instanceType := astSelectorExpr(instance.Package, instance.Name)

		stmts := make([]ast.Stmt, 0)
		// [code] {{Name}} := &{{Type}}{}
		stmts = append(stmts, astDefineStmt(
			astIdent(instanceVar),
			astDeclareExpr(instanceType),
		))
		for _, field := range instance.Fields {
			fieldInstance := field.Instance
			if fieldInstance == "" {
				fieldInstance = field.Name
				if field.Name == "" {
					fieldInstance = utils.ShortType(field.Type)
				}
			}

			if field.IsInject {
				// [code] {{Instance}}.{{FieldName}} = container.{{FieldInstance}}
				stmts = append(stmts, astAssignStmt(
					astSelectorExpr(instanceVar, fieldInstance),
					astSelectorExpr(recvVar, field.Name),
				))
			} else {
				// [code] {{Instance}}.{{FieldName}} = {{FieldInstance}}
				stmts = append(stmts, astAssignStmt(
					astSelectorExpr(instanceVar, fieldInstance),
					astIdent(field.Name),
				))
			}
		}
		// [code] return {{Instance}}
		stmts = append(stmts, &ast.ReturnStmt{
			Results: []ast.Expr{
				astIdent(instanceVar),
			},
		})

		// [code] func (container *Container) {{Constructor}}(
		constructor := instance.Constructor
		if constructor == "" {
			constructor = "New" + provideInstance
		}
		funcDecl := astFuncDecl(
			[]*ast.Field{
				astField(recvVar, astStarExpr(astIdent(global.StructName))),
			},
			constructor,
			param,
			[]*ast.Field{
				{
					Type: astStarExpr(instanceType),
				},
			},
			stmts,
		)

		addDecl(astFile, funcDecl)
	}

}
