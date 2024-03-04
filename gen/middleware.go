package gen

import (
	"github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

func genMiddlewareAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	recvVar := utils.FirstToLower(global.StructName)
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
								global.GenPackage,
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
				switch paramInfo.Source {
				case "ctx":
					args = append(args, astIdent("ctx"))
					break
				case "webCtx":
					args = append(args, astIdent("webCtx"))
					break
				case "inject":
					args = append(args, astSelectorExpr("ctx", paramInfo.Instance))
					break
				default:
					args = append(args, astIdent(paramInfo.Instance))
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

			addDecl(astFile, astFuncDecl(
				[]*ast.Field{
					astField(recvVar, astStarExpr(astIdent(global.StructName))),
				},
				instance.Proxy,
				params,
				[]*ast.Field{
					astField("err", astIdent("error")),
				},
				stmts,
			))
		}
	}
}
