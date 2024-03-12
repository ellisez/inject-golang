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
	for _, webApp := range ctx.SingletonInstances {
		if webInstance, ok := webApp.(*model.WebInstance); ok {
			for _, instance := range webInstance.Middlewares {

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

				instanceCallExpr := astInstanceCallExpr(astSelectorExpr(instance.Package, instance.FuncName), webInstance.Func, ctx, ctxVar)
				stmts = append(stmts, &ast.ReturnStmt{
					Results: []ast.Expr{
						instanceCallExpr,
					},
				})

				genDoc := &ast.Comment{
					Text: fmt.Sprintf("// Generate by annotations from %s.%s", instance.Package, instance.FuncName),
				}
				funcDecl := astInstanceProxyFunc(instance.Func, instance.Instance,
					astField("webCtx",
						astStarExpr(astSelectorExpr("fiber", "Ctx"))))
				funcDecl.Type.Results = &ast.FieldList{List: []*ast.Field{
					astField("err", ast.NewIdent("error")),
				}}
				funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
					{
						Text: fmt.Sprintf("// %s", instance.Proxy),
					},
					genDoc,
				}}
				addDecl(astFile, funcDecl)
				ctx.Methods = append(ctx.Methods, funcDecl)
			}
		}
	}

}
