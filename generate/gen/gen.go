package gen

import (
	"errors"
	"fmt"
	. "github.com/ellisez/inject-golang/generate/global"
	"github.com/ellisez/inject-golang/generate/model"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
)

func DoGen(annotateInfo *model.AnnotateInfo) error {
	if len(annotateInfo.SingletonInstances) == 0 &&
		len(annotateInfo.MultipleInstances) == 0 &&
		len(annotateInfo.FuncInstances) == 0 &&
		len(annotateInfo.MethodInstances) == 0 {
		return nil
	}
	genDir := filepath.Join(RootDirectory, GenPackage)

	err := createGenRootDirectory(genDir)
	if err != nil {
		return err
	}

	err = genContainerFile(annotateInfo, genDir)
	if err != nil {
		return err
	}

	err = genConstructFile(annotateInfo, genDir)
	if err != nil {
		return err
	}

	err = genFuncFile(annotateInfo, genDir)
	if err != nil {
		return err
	}

	err = genMethodFile(annotateInfo, genDir)
	if err != nil {
		return err
	}

	return nil
}

func createGenRootDirectory(genDir string) error {
	stat, err := os.Stat(genDir)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(genDir)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if !stat.IsDir() {
		return errors.New(fmt.Sprintf("%s is not Directory!\ncan not generate code.", genDir))
	}
	return nil
}

func astIdent(text string) *ast.Ident {
	ident := new(ast.Ident)
	ident.Name = text
	return ident
}

func astSelectorExpr(x string, sel string) *ast.SelectorExpr {
	selectorExpr := new(ast.SelectorExpr)
	selectorExpr.X = astIdent(x)
	selectorExpr.Sel = astIdent(sel)
	return selectorExpr
}

func astSelectorExpr1(x ast.Expr, sel string) *ast.SelectorExpr {
	selectorExpr := new(ast.SelectorExpr)
	selectorExpr.X = x
	selectorExpr.Sel = astIdent(sel)
	return selectorExpr
}

func astStarExpr(x ast.Expr) *ast.StarExpr {
	starExpr := new(ast.StarExpr)
	starExpr.X = x
	return starExpr
}

func astDeclareExpr(typeExpr ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.CompositeLit{
			Type: typeExpr,
		},
	}
}

func astDefineStmt(lhs ast.Expr, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			lhs,
		},
		Rhs: []ast.Expr{
			rhs,
		},
	}
}

func astField(name string, typeExpr ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{
			astIdent(name),
		},
		Type: typeExpr,
	}
}

func astAssignStmt(lhs ast.Expr, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			lhs,
		},
		Rhs: []ast.Expr{
			rhs,
		},
	}
}

func astFuncDecl(recv []*ast.Field, name string, params []*ast.Field, results []*ast.Field,
	body []ast.Stmt) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: recv,
		},
		Name: astIdent(name),
		Body: &ast.BlockStmt{
			List: body,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: params,
			},
			Results: &ast.FieldList{
				List: results,
			},
		},
	}
}

func astStructDecl(name string, fields []*ast.Field) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: astIdent(name),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fields,
					},
				},
			},
		},
	}
}

func astImport(astFile *ast.File, importName string, importPath string) *ast.ImportSpec {
	var astImport *ast.ImportSpec
	for _, aImport := range astFile.Imports {
		if importPath == aImport.Path.Value {
			astImport = aImport
			break
		}
	}
	if astImport == nil {
		astImport = &ast.ImportSpec{
			Name: astIdent(importName),
			Path: &ast.BasicLit{
				Value: importPath,
			},
		}

		addImport(astFile, astImport)
	}
	return astImport
}

func addImport(astFile *ast.File, astImport *ast.ImportSpec) {
	if astFile.Imports == nil {
		astFile.Imports = make([]*ast.ImportSpec, 0)
	}
	astFile.Imports = append(astFile.Imports, astImport)

	genDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			astImport,
		},
	}

	addDecl(astFile, genDecl)
}

func addDecl(astFile *ast.File, genDecl ast.Decl) {
	if astFile.Decls == nil {
		astFile.Decls = make([]ast.Decl, 0)
	}
	astFile.Decls = append(astFile.Decls, genDecl)
}
