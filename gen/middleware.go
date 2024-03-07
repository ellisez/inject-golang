package gen

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

func genMiddlewareAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	recvVar := utils.FirstToLower(StructName)
	for _, webApp := range moduleInfo.WebAppInstances {
		for _, instance := range webApp.Middlewares {
			params := []*ast.Field{
				astField("webCtx",
					astStarExpr(astSelectorExpr("fiber", "Ctx"))),
			}

			for _, paramInfo := range instance.Params {
				if paramInfo.Source == "" {
					// [code] {{ParamInstance}} {{ParamType}},
					paramInstance := paramInfo.Instance
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
			stmts = append(stmts, genWebBodyParam(
				instance.BodyParam, instance.PackageInfo, instance.FuncName)...)
			stmts = append(stmts, genWebHeaderParams(
				instance.HeaderParams, instance.PackageInfo, instance.FuncName)...)
			stmts = append(stmts, genWebQueryParams(
				instance.QueryParams, instance.PackageInfo, instance.FuncName)...)
			stmts = append(stmts, genWebPathParams(
				instance.PathParams, instance.PackageInfo, instance.FuncName)...)
			stmts = append(stmts, genWebFormParams(
				instance.FormParams, instance.PackageInfo, instance.FuncName)...)

			args := make([]ast.Expr, 0)
			for _, paramInfo := range instance.Params {
				paramInstance := paramInfo.Instance
				switch paramInfo.Source {
				case "ctx":
					args = append(args, astIdent("ctx"))
					break
				case "webCtx":
					args = append(args, astIdent("webCtx"))
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
					args = append(args, astIdent(paramInstance))
				}

			}

			var fun ast.Expr
			if instance.Recv == nil {
				fun = astSelectorExpr(instance.Package, instance.FuncName)
			} else {
				if instance.Recv.Source == "inject" {
					fun = astSelectorExprRecur(
						astSelectorExpr("ctx", instance.Recv.Instance),
						instance.FuncName,
					)
				} else {
					fun = astSelectorExpr(instance.Recv.Instance, instance.FuncName)
				}
			}
			stmts = append(stmts, &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun:  fun,
						Args: args,
					},
				},
			})

			genDoc := &ast.Comment{
				Text: fmt.Sprintf("// Generate by annotations from %s.%s", instance.Package, instance.FuncName),
			}
			results := []*ast.Field{
				astField("err", astIdent("error")),
			}
			funcDecl := astFuncDecl(
				[]*ast.Field{
					astField(recvVar, astStarExpr(astIdent(StructName))),
				},
				instance.Proxy,
				params,
				results,
				stmts,
			)
			funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
				{
					Text: fmt.Sprintf("// %s", instance.Proxy),
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
}
