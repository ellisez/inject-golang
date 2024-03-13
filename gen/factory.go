package gen

import (
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"path"
	"path/filepath"
)

func genFactoryFile(ctx *model.Ctx, dir string) error {
	fileDir := filepath.Join(dir, GenFactoryPackage)
	filename := filepath.Join(fileDir, GenFactoryFilename)

	astFile := &ast.File{
		Name:  ast.NewIdent(GenFactoryPackage),
		Scope: ast.NewScope(nil),
	}

	genFactoryImportAst(ctx, astFile)

	genFactoryNewAst(ctx, astFile)

	astFile, err := utils.OptimizeCode(filename, astFile, ctx,
		`// Code generated by "inject-golang"; DO NOT EDIT.`)
	if err != nil {
		return err
	}

	return utils.GenerateCode(filename, astFile, ctx)
}

func genFactoryImportAst(ctx *model.Ctx, astFile *ast.File) {
	addImport(astFile, ctx, "", path.Join(Mod.Package, GenPackage, "internal"))
}

func genFactoryNewAst(ctx *model.Ctx, astFile *ast.File) {
	funcDecl := astFuncDecl(
		nil,
		"New",
		nil,
		[]*ast.Field{
			astField("", astStarExpr(astSelectorExpr(GenInternalPackage, CtxType))),
		},
		[]ast.Stmt{
			&ast.ReturnStmt{Results: []ast.Expr{
				&ast.CallExpr{
					Fun: astSelectorExpr(GenInternalPackage, "New"),
				},
			}},
		},
	)
	if ctx.Doc != nil {
		funcDecl.Doc = &ast.CommentGroup{List: ctx.Doc}
	}
	addDecl(astFile, funcDecl)
}
