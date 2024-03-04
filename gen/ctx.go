package gen

import (
	"github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// gen_ctx.go
func genCtxFile(moduleInfo *model.ModuleInfo, dir string) error {
	filename := filepath.Join(dir, global.GenCtxFilename)

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	astFile := &ast.File{
		Name:  astIdent(global.GenPackage),
		Scope: ast.NewScope(nil),
	}

	genCtxImportsAst(moduleInfo, astFile)

	genCtxStructAst(moduleInfo, astFile)

	genCtxConstructorAst(moduleInfo, astFile)

	err = format.Node(file, token.NewFileSet(), astFile)
	if err != nil {
		return err
	}
	return nil
}

func genCtxImportsAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {

	for _, instance := range moduleInfo.SingletonInstances {
		astImport(astFile, "", instance.Import)
		if instance.Imports != nil {
			for _, importInfo := range instance.Imports {
				importName := importInfo.Name
				if importName == "_" {
					importName = ""
				}
				astImport(astFile, importName, importInfo.Path)
			}
		}
	}
	if moduleInfo.WebAppInstances != nil {
		astImport(astFile, "", "github.com/gofiber/fiber/v2")
	}
	addImportDecl(astFile)

}

// # gen segment: Struct #
func genCtxStructAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	fields := make([]*ast.Field, 0)
	for _, instance := range moduleInfo.SingletonInstances {
		fieldName := instance.Instance

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
	if moduleInfo.WebAppInstances != nil {
		for _, instance := range moduleInfo.WebAppInstances {
			fields = append(fields, astField(
				instance.WebApp,
				astStarExpr(
					astSelectorExpr(
						"fiber",
						"App",
					),
				),
			))
		}
	}

	genDecl := astStructDecl(
		global.StructName,
		fields,
	)

	addDecl(astFile, genDecl)

}

// # gen segment: Singleton instance #
func genCtxConstructorAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	varName := "ctx"
	stmts := make([]ast.Stmt, 0)
	// [code] ctx := &ProvideContainer{}
	stmts = append(stmts, astDefineStmt(
		astIdent(varName),
		astDeclareRef(astIdent(global.StructName), nil),
	))

	assignStmts := make([]ast.Stmt, 0)
	postStmts := make([]ast.Stmt, 0)
	for _, instance := range moduleInfo.SingletonInstances {
		provideInstance := instance.Instance

		if instance.PreConstruct != "" {
			// [code] ctx.{{Instance}} = {{PreConstruct}}()
			var caller ast.Expr
			if !strings.Contains(instance.PreConstruct, ".") {
				if moduleInfo.HasFunc(instance.PreConstruct) {
					caller = astSelectorExpr(varName, instance.PreConstruct)
				} else {
					utils.Failuref("@preConstruct %s, No matching function, try to specify Package Name, at %s{}", instance.PreConstruct, instance.Name)
				}
			} else {
				caller = utils.TypeToAst(instance.PreConstruct)
			}
			stmts = append(stmts, astAssignStmt(
				astSelectorExpr(varName, provideInstance),
				&ast.CallExpr{
					Fun: caller,
				},
			))
		} else {
			// [code] ctx.{{Instance}} = &{{Package}}.{{Name}}{}
			stmts = append(stmts, astAssignStmt(
				astSelectorExpr(varName, provideInstance),
				astDeclareRef(
					astSelectorExpr(
						instance.Package,
						instance.Name,
					),
					nil,
				),
			))
		}

		for _, field := range instance.Fields {
			if field.Source == "inject" {
				fieldInstance := field.Instance
				if fieldInstance == "Ctx" {
					// [code] ctx.{{Instance}}.{{FieldInstance}} = ctx
					assignStmts = append(assignStmts, astAssignStmt(
						astSelectorExprRecur(astSelectorExpr(varName, provideInstance), fieldInstance),
						astIdent(varName),
					))
				} else {
					if fieldInstance == "Ctx" {
						assignStmts = append(assignStmts, astAssignStmt(
							astSelectorExprRecur(astSelectorExpr(varName, provideInstance), fieldInstance),
							astIdent(varName),
						))
					} else {
						// [code] ctx.{{Instance}}.{{FieldInstance}} = ctx.{{StructInstance}}
						if !moduleInfo.HasInstance(fieldInstance) {
							utils.Failuref("%s, \"%s\" No matching Instance, at %s{}", field.Comment, fieldInstance, instance.Name)
						}
						assignStmts = append(assignStmts, astAssignStmt(
							astSelectorExprRecur(astSelectorExpr(varName, provideInstance), fieldInstance),
							astSelectorExpr(varName, fieldInstance),
						))
					}
				}
			}
		}

		if instance.PostConstruct != "" {
			// [code] {{PostConstruct}}(ctx.{{Instance}})
			var caller ast.Expr
			if !strings.Contains(instance.PostConstruct, ".") {
				if moduleInfo.HasFunc(instance.PostConstruct) {
					caller = astSelectorExpr(varName, instance.PostConstruct)
				} else {
					utils.Failuref("@postConstruct %s, No matching function, try to specify Package Name, at %s{}", instance.PreConstruct, instance.Name)
				}
			} else {
				caller = utils.TypeToAst(instance.PostConstruct)
			}
			postStmts = append(postStmts, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: caller,
					Args: []ast.Expr{
						astSelectorExpr(varName, provideInstance),
					},
				},
			})
		}
	}

	for _, instance := range moduleInfo.WebAppInstances {
		// [code] ctx.{{WebApp}} = fiber.New()
		stmts = append(stmts, astAssignStmt(
			astSelectorExpr(varName, instance.WebApp),
			&ast.CallExpr{
				Fun: astSelectorExpr(
					"fiber",
					"New",
				),
			},
		))
	}

	stmts = append(stmts, assignStmts...)
	stmts = append(stmts, postStmts...)
	// [code] return ctx
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
				Type: astStarExpr(astIdent(global.StructName)),
			},
		},
		stmts,
	)

	addDecl(astFile, funcDecl)
}
