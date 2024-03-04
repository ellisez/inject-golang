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

// gen_constructor.go
func genConstructorFile(moduleInfo *model.ModuleInfo, dir string) error {
	filename := filepath.Join(dir, global.GenConstructorFilename)

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	astFile := &ast.File{
		Name:  astIdent(global.GenPackage),
		Scope: ast.NewScope(nil),
	}

	genConstructorImportsAst(moduleInfo, astFile)

	genConstructorAst(moduleInfo, astFile)

	err = format.Node(file, token.NewFileSet(), astFile)
	if err != nil {
		return err
	}
	return nil
}

func genConstructorImportsAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {

	for _, instance := range moduleInfo.MultipleInstances {
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
	addImportDecl(astFile)
}

// # gen segment: Multiple instance #
func genConstructorAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {

	for _, instance := range moduleInfo.MultipleInstances {
		recvVar := utils.FirstToLower(global.StructName)
		params := make([]*ast.Field, 0)
		for _, field := range instance.Fields {
			if field.Source == "" {
				// [code] {{FieldInstance}} {{FieldType}},
				fieldInstance := field.Instance

				params = append(params,
					astField(
						fieldInstance,
						utils.AccessType(
							field.Type,
							instance.Package,
							global.GenPackage,
						),
					),
				)
			}
		}

		provideInstance := instance.Instance
		if provideInstance == "" || provideInstance == "_" {
			provideInstance = instance.Name
		}

		instanceVar := utils.FirstToLower(provideInstance)
		instanceType := astSelectorExpr(instance.Package, instance.Name)

		stmts := make([]ast.Stmt, 0)
		if instance.PreConstruct != "" {
			// [code] {{Instance}} := {{PreConstruct}}()
			var caller ast.Expr
			if !strings.Contains(instance.PreConstruct, ".") {
				if moduleInfo.HasFunc(instance.PreConstruct) {
					caller = astSelectorExpr(recvVar, instance.PreConstruct)
				} else {
					utils.Failuref("@preConstruct %s, No matching function, must be to specify Package Name, at %s{}", instance.PreConstruct, instance.Name)
				}
			} else {
				caller = utils.TypeToAst(instance.PreConstruct)
			}
			stmts = append(stmts, astDefineStmt(
				astIdent(instanceVar),
				&ast.CallExpr{
					Fun: caller,
				},
			))
		} else {
			// [code] {{Instance}} := &{{Package}}.{{Name}}{}
			stmts = append(stmts, astDefineStmt(
				astIdent(instanceVar),
				astDeclareRef(instanceType, nil),
			))
		}
		for _, field := range instance.Fields {
			fieldInstance := field.Instance

			if field.Source == "inject" {
				if fieldInstance == "Ctx" {
					// [code] {{Instance}}.{{FieldName}} = ctx
					stmts = append(stmts, astAssignStmt(
						astSelectorExpr(instanceVar, fieldInstance),
						astIdent(recvVar),
					))
				} else {
					// [code] {{Instance}}.{{FieldName}} = ctx.{{FieldInstance}}
					if !moduleInfo.HasInstance(fieldInstance) {
						utils.Failuref("%s, \"%s\" No matching Instance, at %s{}", field.Comment, fieldInstance, instance.Name)
					}
					stmts = append(stmts, astAssignStmt(
						astSelectorExpr(instanceVar, fieldInstance),
						astSelectorExpr(recvVar, field.Name),
					))
				}
			} else {
				// [code] {{Instance}}.{{FieldName}} = {{FieldInstance}}
				stmts = append(stmts, astAssignStmt(
					astSelectorExpr(instanceVar, fieldInstance),
					astIdent(field.Name),
				))
			}
		}

		if instance.PostConstruct != "" {
			// [code] {{PostConstruct}}({{Instance}})
			var caller ast.Expr
			if !strings.Contains(instance.PostConstruct, ".") {
				if moduleInfo.HasFunc(instance.PostConstruct) {
					caller = astSelectorExpr(recvVar, instance.PostConstruct)
				} else {
					utils.Failuref("@postConstruct %s, No matching function, must be to specify Package Name, at %s{}", instance.PreConstruct, instance.Name)
				}
			} else {
				caller = utils.TypeToAst(instance.PostConstruct)
			}
			stmts = append(stmts, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: caller,
					Args: []ast.Expr{
						astIdent(instanceVar),
					},
				},
			})
		}
		// [code] return {{Instance}}
		stmts = append(stmts, &ast.ReturnStmt{
			Results: []ast.Expr{
				astIdent(instanceVar),
			},
		})

		// [code] func (ctx *Ctx) New{{Instance}}(
		constructor := "New" + provideInstance

		funcDecl := astFuncDecl(
			[]*ast.Field{
				astField(recvVar, astStarExpr(astIdent(global.StructName))),
			},
			constructor,
			params,
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
