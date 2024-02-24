package gen

import (
	"github.com/ellisez/inject-golang/generate/global"
	"github.com/ellisez/inject-golang/generate/model"
	"github.com/ellisez/inject-golang/generate/utils"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
)

// __gen_func.go
func genFuncFile(annotateInfo *model.AnnotateInfo, dir string) error {
	filename := filepath.Join(dir, "__gen_container.go")
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	astFile := &ast.File{
		Name:  astIdent(global.GenPackage),
		Scope: ast.NewScope(nil),
	}

	genFuncImportsAst(annotateInfo, astFile)

	genFuncAst(annotateInfo, astFile)

	err = format.Node(file, token.NewFileSet(), astFile)
	if err != nil {
		return err
	}
	return nil
}
func genFuncImportsAst(annotateInfo *model.AnnotateInfo, astFile *ast.File) {

	for _, instance := range annotateInfo.FuncInstances {
		astImport(astFile, "", instance.Dirname)
	}
}

// # gen segment: Func inject #
func genFuncAst(annotateInfo *model.AnnotateInfo, astFile *ast.File) {
	for _, instance := range annotateInfo.FuncInstances {
		recvVar := utils.FirstToLower(global.StructName)
		param := make([]*ast.Field, 0)
		for _, paramInfo := range instance.NormalParams {
			// [code] {{ParamInstance}} {{ParamType}},
			paramInstance := paramInfo.Instance
			if paramInstance == "" {
				paramInstance = paramInfo.Name
				if paramInfo.Name == "" {
					paramInstance = utils.ShortType(paramInfo.Type)
				}
			}
			param = append(param, astField(paramInstance, utils.TypeToAst(paramInfo.Type)))
		}

		stmts := make([]ast.Stmt, 0)
		args := make([]ast.Expr, 0)
		for _, paramInfo := range instance.Params {
			paramInstance := paramInfo.Instance
			if paramInstance == "" {
				paramInstance = paramInfo.Name
				if paramInfo.Name == "" {
					paramInstance = utils.ShortType(paramInfo.Type)
				}
			}

			if paramInfo.IsInject {
				// [code] container.{{ParamInstance}},
				args = append(args, astSelectorExpr(recvVar, paramInstance))
			} else {
				// [code] {{ParamInstance}},
				args = append(args, astIdent(paramInstance))
			}
		}
		stmts = append(stmts, &ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.CallExpr{
					Fun:  astSelectorExpr(instance.Package, instance.FuncName),
					Args: args,
				},
			},
		})

		results := make([]*ast.Field, 0)
		for _, result := range instance.Results {
			results = append(results, astField(result.Name, utils.TypeToAst(result.Type)))
		}

		// [code] func (container *Container) {{Proxy}}(
		funcDecl := astFuncDecl(
			[]*ast.Field{
				astField(recvVar, astStarExpr(astIdent(global.StructName))),
			},
			instance.Proxy,
			param,
			results,
			stmts,
		)

		addDecl(astFile, funcDecl)
	}

}
