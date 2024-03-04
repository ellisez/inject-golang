package gen

import (
	"fmt"
	"github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/token"
	"path/filepath"
)

func DoGen(moduleInfo *model.ModuleInfo) error {
	if len(moduleInfo.SingletonInstances) == 0 &&
		len(moduleInfo.MultipleInstances) == 0 &&
		len(moduleInfo.FuncInstances) == 0 &&
		len(moduleInfo.MethodInstances) == 0 {
		return nil
	}
	genDir := filepath.Join(moduleInfo.Dirname, global.GenPackage)

	err := utils.CreateDirectoryIfNotExists(genDir)
	if err != nil {
		return err
	}

	err = genCtxFile(moduleInfo, genDir)
	if err != nil {
		return err
	}

	err = genConstructorFile(moduleInfo, genDir)
	if err != nil {
		return err
	}

	err = genFuncFile(moduleInfo, genDir)
	if err != nil {
		return err
	}

	err = genMethodFile(moduleInfo, genDir)
	if err != nil {
		return err
	}

	err = genWebFile(moduleInfo, genDir)
	if err != nil {
		return err
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

func astSelectorExprRecur(x ast.Expr, sel string) *ast.SelectorExpr {
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

func astDeclareRef(typeExpr ast.Expr, elts []ast.Expr) *ast.UnaryExpr {
	starExpr, ok := typeExpr.(*ast.StarExpr)
	if ok {
		typeExpr = starExpr.X
	}
	return &ast.UnaryExpr{
		Op: token.AND,
		X:  astDeclareExpr(typeExpr, elts),
	}
}

func astDeclareExpr(typeExpr ast.Expr, elts []ast.Expr) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: typeExpr,
		Elts: elts,
	}
}

func astIntExpr(number string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: number,
	}
}

func astStringExpr(text string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: fmt.Sprintf(`"%s"`, text),
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

func astDefineStmtMany(lhs []ast.Expr, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: lhs,
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

func astAssignStmtMany(lhs []ast.Expr, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: lhs,
		Rhs: []ast.Expr{
			rhs,
		},
	}
}

func astFuncDecl(recv []*ast.Field, name string, params []*ast.Field, results []*ast.Field,
	body []ast.Stmt) *ast.FuncDecl {
	var recvAst *ast.FieldList
	if recv != nil {
		recvAst = &ast.FieldList{
			List: recv,
		}
	}
	return &ast.FuncDecl{
		Recv: recvAst,
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
	importPath = fmt.Sprintf(`"%s"`, importPath)
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

		addImportSpec(astFile, astImport)
	}
	return astImport
}

func astTypeToDeclare(typeExpr ast.Expr) ast.Expr {
	switch typeExpr.(type) {
	case *ast.StarExpr:
		return astDeclareRef(typeExpr.(*ast.StarExpr).X, nil)
	case *ast.Ident:
		if utils.IsFirstLower(typeExpr.(*ast.Ident).String()) {
			return nil
		}
	}
	return typeExpr
}

func addImportSpec(astFile *ast.File, astImport *ast.ImportSpec) {
	if astFile.Imports == nil {
		astFile.Imports = make([]*ast.ImportSpec, 0)
	}
	astFile.Imports = append(astFile.Imports, astImport)
}

func addImportDecl(astFile *ast.File) {
	specs := make([]ast.Spec, len(astFile.Imports))
	for i, spec := range astFile.Imports {
		specs[i] = spec
	}
	genDecl := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: specs,
	}

	addDecl(astFile, genDecl)
}

func addDecl(astFile *ast.File, genDecl ast.Decl) {
	if astFile.Decls == nil {
		astFile.Decls = make([]ast.Decl, 0)
	}
	astFile.Decls = append(astFile.Decls, genDecl)
}
