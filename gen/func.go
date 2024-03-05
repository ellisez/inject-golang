package gen

import (
	"fmt"
	"github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
)

// gen_func.go
func genFuncFile(moduleInfo *model.ModuleInfo, dir string) error {
	filename := filepath.Join(dir, global.GenFuncFilename)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	astFile := &ast.File{
		Name:  astIdent(global.GenPackage),
		Scope: ast.NewScope(nil),
	}

	genFuncImportsAst(moduleInfo, astFile)

	genFuncAst(moduleInfo, astFile)

	addFileDoc(astFile, "\n// Code generate by \"inject-golang func\"; DO NOT EDIT.")

	err = format.Node(file, token.NewFileSet(), astFile)
	if err != nil {
		return err
	}
	return nil
}
func genFuncImportsAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {

	for _, instance := range moduleInfo.FuncInstances {
		astImport(astFile, "", instance.Import)
		for _, importInfo := range instance.Imports {
			importName := importInfo.Name
			if importName == "_" {
				importName = ""
			}
			astImport(astFile, importName, importInfo.Path)
		}
	}
	addImportDecl(astFile)
}

// # gen segment: Func inject #
func genFuncAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	recvVar := utils.FirstToLower(global.StructName)

	for _, instance := range moduleInfo.FuncInstances {
		params := make([]*ast.Field, 0)
		for _, paramInfo := range instance.Params {
			if paramInfo.Source == "" {
				// [code] {{ParamInstance}} {{ParamType}},
				paramInstance := paramInfo.Instance
				if paramInstance == "" {
					paramInstance = paramInfo.Name
					if paramInfo.Name == "" {
						paramInstance = utils.TypeShortName(paramInfo.Type)
					}
				}
				params = append(params,
					astField(paramInstance,
						utils.AccessType(
							paramInfo.Type,
							instance.Package,
							global.GenPackage,
						),
					),
				)
			}
		}

		stmts := make([]ast.Stmt, 0)
		args := make([]ast.Expr, 0)
		for _, paramInfo := range instance.Params {
			paramInstance := paramInfo.Instance

			switch paramInfo.Source {
			case "ctx":
				// [code] ctx,
				args = append(args, astIdent(recvVar))
				break
			case "inject":
				// [code] ctx.{{ParamInstance}},
				if !moduleInfo.HasInstance(paramInstance) {
					utils.Failuref("%s, \"%s\" No matching Instance, at %s()", paramInfo.Comment, paramInstance, instance.FuncName)
				}
				args = append(args, astSelectorExpr(recvVar, paramInstance))

			default:
				// [code] {{ParamInstance}},
				args = append(args, astIdent(paramInstance))
			}
		}
		if len(instance.Results) == 0 {
			stmts = append(stmts, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun:  astSelectorExpr(instance.Package, instance.FuncName),
					Args: args,
				},
			})
		} else {
			stmts = append(stmts, &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun:  astSelectorExpr(instance.Package, instance.FuncName),
						Args: args,
					},
				},
			})
		}

		results := make([]*ast.Field, 0)
		for _, result := range instance.Results {
			results = append(results, astField(result.Name, result.Type))
		}

		// [code] func (ctx *Container) {{Proxy}}(
		funcDecl := astFuncDecl(
			[]*ast.Field{
				astField(recvVar, astStarExpr(astIdent(global.StructName))),
			},
			instance.Proxy,
			params,
			results,
			stmts,
		)
		funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
			{
				Text: fmt.Sprintf("\n// %s", instance.Proxy),
			},
			{
				Text: fmt.Sprintf("// Generate by annotations from %s.%s", instance.Package, instance.FuncName),
			},
		}}

		addDecl(astFile, funcDecl)
	}

}
