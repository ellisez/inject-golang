package gen

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"os"
	"path/filepath"
)

// gen_func.go
func genFuncFile(moduleInfo *model.ModuleInfo, dir string) error {
	fileDir := filepath.Join(dir, GenInternalPackage)
	filename := filepath.Join(fileDir, GenFuncFilename)

	if moduleInfo.FuncInstances == nil {
		err := os.Remove(filename)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		return nil
	}

	astFile := &ast.File{
		Name:  astIdent(GenInternalPackage),
		Scope: ast.NewScope(nil),
	}

	genFuncImportsAst(moduleInfo, astFile)

	genFuncAst(moduleInfo, astFile)

	astFile, err := utils.FixErrors(filename, astFile, moduleInfo,
		"// Code generated by \"inject-golang -m func\"; DO NOT EDIT.")
	if err != nil {
		return err
	}

	return utils.GenerateCode(filename, astFile, moduleInfo)
}
func genFuncImportsAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {

	for _, instance := range moduleInfo.FuncInstances {
		addImport(astFile, moduleInfo, "", instance.Import)
		for _, importInfo := range instance.Imports {
			importName := importInfo.Name
			if importName == "_" {
				importName = ""
			}
			addImport(astFile, moduleInfo, importName, importInfo.Path)
		}
	}
}

// # gen segment: Func inject #
func genFuncAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	recvVar := utils.FirstToLower(StructName)

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
							GenPackage,
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
				injectMode := ""
				if moduleInfo.HasSingleton(paramInstance) {
					// [code] ctx.{{ParamInstance}}(),
					args = append(args,
						&ast.CallExpr{
							Fun: astSelectorExpr(recvVar, paramInstance),
						})
					injectMode = "singleton"
				}

				if injectMode == "" {
					if moduleInfo.HasWebApp(paramInstance) {
						// [code] ctx.{{ParamInstance}}(),
						args = append(args,
							&ast.CallExpr{
								Fun: astSelectorExpr(recvVar, paramInstance),
							})
						injectMode = "webApp"
					}
				}

				if injectMode == "" {
					if moduleInfo.HasMultiple(paramInstance) {
						// [code] ctx.New{{ParamInstance}}(),
						args = append(args, &ast.CallExpr{
							Fun: astSelectorExpr(recvVar, "New"+paramInstance),
						})
						injectMode = "multiple"
					}
				}

				if injectMode == "" {
					utils.Failuref("%s, \"%s\" No matching Instance, at %s()", paramInfo.Comment, paramInstance, instance.FuncName)
				}
				break
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
				astField(recvVar, astStarExpr(astIdent(StructName))),
			},
			instance.Proxy,
			params,
			results,
			stmts,
		)
		genDoc := &ast.Comment{
			Text: fmt.Sprintf("// Generate by annotations from %s.%s", instance.Package, instance.FuncName),
		}
		funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
			{
				Text: fmt.Sprintf("\n// %s", instance.Proxy),
			},
			genDoc,
		}}
		addDecl(astFile, funcDecl)

		/// interface method field
		methodField := astField(
			instance.Proxy,
			&ast.FuncType{
				Params:  &ast.FieldList{List: params},
				Results: &ast.FieldList{List: results},
			},
		)
		methodField.Comment = &ast.CommentGroup{List: []*ast.Comment{
			genDoc,
		}}
		moduleInfo.CtxMethodFields = append(moduleInfo.CtxMethodFields, methodField)
	}

}
