package utils

import (
	"bytes"
	"github.com/ellisez/inject-golang/model"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func FixErrors(filename string, astFile *ast.File, moduleInfo *model.ModuleInfo, doc string) (*ast.File, error) {
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
	///////////
	buffer := &bytes.Buffer{}
	err := format.Node(buffer, moduleInfo.FileSet, astFile)
	if err != nil {
		return nil, err
	}
	newAstFile, err := parser.ParseFile(moduleInfo.FileSet, filename, buffer, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	/// remove unused imports
	UnusedImports(newAstFile)
	return newAstFile, nil
}

func GenerateCode(filename string, astFile *ast.File, moduleInfo *model.ModuleInfo) error {
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

	err = format.Node(file, moduleInfo.FileSet, astFile)
	if err != nil {
		return err
	}
	return nil
}
