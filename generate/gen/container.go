package gen

import (
	"github.com/ellisez/inject-golang/generate/global"
	"github.com/ellisez/inject-golang/generate/model"
	"github.com/ellisez/inject-golang/generate/utils"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
)

// __gen_container.go
func genContainerFile(annotateInfo *model.AnnotateInfo, dir string) error {
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

	genContainerImportsAst(annotateInfo, astFile)

	genContainerStructAst(annotateInfo, astFile)

	genContainerConstructorAst(annotateInfo, astFile)

	err = format.Node(file, token.NewFileSet(), astFile)
	if err != nil {
		return err
	}
	return nil
}

func genContainerImportsAst(annotateInfo *model.AnnotateInfo, astFile *ast.File) {

	for _, instance := range annotateInfo.SingletonInstances {
		astImport(astFile, "", instance.Dirname)
	}

}

// # gen segment: Struct #
func genContainerStructAst(annotateInfo *model.AnnotateInfo, astFile *ast.File) {
	fields := make([]*ast.Field, 0)
	for _, instance := range annotateInfo.SingletonInstances {
		fieldName := instance.Instance
		if fieldName == "" {
			fieldName = instance.Name
		}

		fields = append(fields, astField(
			fieldName,
			astStarExpr(
				astSelectorExpr(
					instance.Package,
					instance.Name,
				),
			),
		))
	}

	genDecl := astStructDecl(
		global.StructName,
		fields,
	)

	addDecl(astFile, genDecl)

}

// # gen segment: Singleton instance #
func genContainerConstructorAst(annotateInfo *model.AnnotateInfo, astFile *ast.File) {
	varName := "container"
	stmts := make([]ast.Stmt, 0)
	// [code] container := &ProvideContainer{}
	stmts = append(stmts, astDefineStmt(
		astIdent(varName),
		astDeclareExpr(astIdent(global.StructName)),
	))

	assignStmts := make([]ast.Stmt, 0)
	for _, instance := range annotateInfo.SingletonInstances {
		// [code] container.{{Instance}} = &{{Package}}.{{Name}}{}
		provideInstance := instance.Instance
		if provideInstance == "" {
			provideInstance = instance.Name
		}

		stmts = append(stmts, astAssignStmt(
			astSelectorExpr(varName, provideInstance),
			astDeclareExpr(
				astSelectorExpr(
					instance.Package,
					instance.Name,
				),
			),
		))

		for _, field := range instance.InjectFields {
			// [code] container.{{Instance}}.{{FieldInstance}} = container.{{StructInstance}}
			fieldInstance := field.Instance
			if fieldInstance == "" {
				fieldInstance = field.Name
				if field.Name == "" {
					fieldInstance = utils.ShortType(field.Type)
				}
			}

			assignStmts = append(assignStmts, astAssignStmt(
				astSelectorExpr1(astSelectorExpr(varName, provideInstance), fieldInstance),
				astSelectorExpr(varName, fieldInstance),
			))
		}
	}

	stmts = append(stmts, assignStmts...)
	// [code] return container
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			astIdent(varName),
		},
	})

	funcDecl := astFuncDecl(
		nil,
		"New",
		nil,
		[]*ast.Field{
			{
				Type: astIdent(global.StructName),
			},
		},
		stmts,
	)

	addDecl(astFile, funcDecl)
}
