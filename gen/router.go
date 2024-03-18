package gen

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

func genRouterAst(ctx *model.Ctx, astFile *ast.File) {
	ctxVar := utils.FirstToLower(CtxType)
	for _, key := range ctx.SingletonInstances.Keys {
		webApp := ctx.SingletonOf(key)
		if webInstance, ok := webApp.(*model.WebInstance); ok {
			for _, instance := range webInstance.Routers {

				stmts := make([]ast.Stmt, 0)
				stmts = append(stmts, genWebBodyParam(
					instance.BodyParam, instance.Package, instance.Func)...)
				stmts = append(stmts, genWebHeaderParams(
					instance.HeaderParams, instance.Package, instance.Func)...)
				stmts = append(stmts, genWebQueryParams(
					instance.QueryParams, instance.Package, instance.Func)...)
				stmts = append(stmts, genWebPathParams(
					instance.PathParams, instance.Package, instance.Func)...)
				stmts = append(stmts, genWebFormParams(
					instance.FormParams, instance.Package, instance.Func)...)

				instanceCallExpr := astInstanceCallExpr(astSelectorExpr(instance.Package, instance.FuncName), instance.Func, ctx, ctxVar)
				stmts = append(stmts, &ast.ReturnStmt{
					Results: []ast.Expr{
						instanceCallExpr,
					},
				})

				funcDecl := astInstanceProxyFunc(instance.Func, instance.Instance,
					astField("webCtx",
						astStarExpr(astSelectorExpr("fiber", "Ctx"))))
				funcDecl.Body = &ast.BlockStmt{List: stmts}
				funcDecl.Type.Results = &ast.FieldList{List: []*ast.Field{
					astField("err", ast.NewIdent("error")),
				}}
				funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{{
					Text: fmt.Sprintf("// Generate by annotations from %s.%s", instance.Package, instance.FuncName),
				}}}
				addDecl(astFile, funcDecl)
				ctx.Methods[funcDecl.Name.String()] = funcDecl
			}
		}
	}

}
