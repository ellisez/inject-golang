package utils

import (
	"bytes"
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func GenerateCode(filename string, astFile *ast.File, ctx *model.Ctx, doc string) error {
	if astFile.Imports != nil {
		specs := make([]ast.Spec, len(astFile.Imports))
		for i, importSpec := range astFile.Imports {
			specs[i] = importSpec
		}

		var docAst *ast.CommentGroup
		if doc != "" {
			docAst = &ast.CommentGroup{List: []*ast.Comment{
				{
					Text: doc,
				},
			}}
		}

		astFile.Decls = append([]ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: specs,
				Doc:   docAst,
			},
		}, astFile.Decls...)
	}

	/// write code
	err := generateCode(filename, astFile, ctx)
	if err != nil {
		return err
	}
	// format code
	buffer := &bytes.Buffer{}
	err = format.Node(buffer, ctx.FileSet, astFile)
	if err != nil {
		return err
	}
	newAstFile, err := parser.ParseFile(ctx.FileSet, filename, buffer, parser.ParseComments)
	if err != nil {
		return err
	}

	/// remove unused imports
	unusedImports(newAstFile)

	if newAstFile.Name.String() == GenPackage {
		/// circular dependency
		circularDependency(newAstFile, ctx)
	}

	// write again
	err = generateCode(filename, newAstFile, ctx)
	if err != nil {
		return err
	}
	return nil
}

func generateCode(filename string, astFile *ast.File, ctx *model.Ctx) error {
	fileDir := filepath.Dir(filename)
	err := CreateDirectoryIfNotExists(fileDir)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	err = format.Node(file, ctx.FileSet, astFile)
	if err != nil {
		return err
	}
	return nil
}

func unusedImports(astFile *ast.File) {
	var imports []*ast.ImportSpec
	addImport := func(spec *ast.ImportSpec) {
		for _, importSpec := range imports {
			if importSpec.Path == spec.Path {
				return
			}
		}
		imports = append(imports, spec)
	}
	for _, spec := range astFile.Imports {
		for _, ident := range astFile.Unresolved {
			specName := spec.Name
			if specName != nil {
				if ident.String() == specName.String() {
					addImport(spec)
				}
			} else {
				specPath := spec.Path.Value
				importPath := specPath[1 : len(specPath)-1]
				defaultPackageName, _ := GetPackageNameFromImport(importPath)
				if defaultPackageName == ident.String() {
					addImport(spec)
				}
			}
		}
	}
	astFile.Imports = imports

	// copy to genDecl
	var specs []ast.Spec
	for _, spec := range imports {
		specs = append(specs, spec)
	}

	// find import genDecl
	var importDecl *ast.GenDecl
	if len(astFile.Decls) > 0 {
		genDecl, ok := astFile.Decls[0].(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT {
			if specs == nil {
				astFile.Decls = astFile.Decls[1:]
			} else {
				genDecl.Specs = specs
			}
			importDecl = genDecl
		}
	}

	// if no import decl
	if importDecl == nil {
		astFile.Decls = append([]ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: specs,
			},
		}, astFile.Decls...)
	}
}

func circularDependency(astFile *ast.File, ctx *model.Ctx) {
	for _, spec := range astFile.Imports {
		for key, ctxFields := range ctx.InjectCtxMap {
			if StringLit(spec.Path) == key {
				var fails []string
				for _, ctxField := range ctxFields {
					fails = append(fails, fmt.Sprintf(`%s %s, "@injectCtx" is not allowed due to circular dependency`, ctxField.Loc.String(), ctxField.Comment))
				}
				Failure(fails...)
			}
		}
	}

}
func GetterOf(instance string) string {
	return FirstToUpper(instance)
}

func SetterOf(instance string) string {
	return "Set" + FirstToUpper(instance)
}
