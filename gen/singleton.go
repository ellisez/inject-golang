package gen

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
)

// gen_ctx.go
func genSingletonFile(ctx *model.Ctx, dir string) error {
	fileDir := filepath.Join(dir, GenInternalPackage)
	filename := filepath.Join(fileDir, GenSingletonFilename)

	if ctx.SingletonInstances == nil {
		err := os.Remove(filename)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		return nil
	}

	astFile := &ast.File{
		Name:  ast.NewIdent(GenInternalPackage),
		Scope: ast.NewScope(nil),
	}

	genSingletonImportsAst(ctx, astFile)

	genSingletonStructAst(ctx, astFile)

	genSingletonGetterAndSetterAst(ctx, astFile)

	genSingletonNewAst(ctx, astFile)

	return utils.GenerateCode(filename, astFile, ctx,
		`// Code generated by "inject-golang -m singleton"; DO NOT EDIT.`)
}

func genSingletonImportsAst(ctx *model.Ctx, astFile *ast.File) {

	for _, instance := range ctx.SingletonInstances {
		if _, ok := instance.(*model.Provide); ok {
			for _, importNode := range instance.GetImports() {
				importName := importNode.Name
				if importName == "_" {
					importName = ""
				}
				addImport(astFile, ctx, importName, importNode.Path)
			}
		}
	}
	if ctx.HasWebInstance {
		addImport(astFile, ctx, "", "github.com/gofiber/fiber/v2")
	}
}

// # gen segment: Struct #
func genSingletonStructAst(ctx *model.Ctx, astFile *ast.File) {
	var fields []*ast.Field
	for _, instance := range ctx.SingletonInstances {
		instanceName := instance.GetInstance()
		instanceType := instance.GetType()

		fieldName := utils.FirstToLower(instanceName)

		fieldType := instanceType

		field := astField(
			fieldName,
			fieldType,
		)
		fields = append(fields, field)
	}

	structDecl := astStructDecl(
		CtxType,
		fields,
	)

	addDecl(astFile, structDecl)
}

func genSingletonGetterAndSetterAst(ctx *model.Ctx, astFile *ast.File) {
	/// Getter / Setter
	for _, instance := range ctx.SingletonInstances {
		instanceName := instance.GetInstance()
		instanceType := instance.GetType()
		instanceFunc := instance.GetFunc()

		fieldName := utils.FirstToLower(instanceName)
		fieldGetter := utils.GetterOf(instanceName)
		fieldSetter := utils.SetterOf(instanceName)
		fieldType := instanceType

		var doc string
		if instanceFunc.Package == "" {
			doc = "// Generate by system"
		} else {
			doc = fmt.Sprintf("// Generate by annotations from %s.%s", instanceFunc.Package, instanceFunc.FuncName)
		}

		getterDecl := astCtxGetter(
			doc,
			fieldGetter,
			fieldName,
			fieldType,
		)
		addDecl(astFile, getterDecl)
		ctx.Methods = append(ctx.Methods, getterDecl)

		setterDecl := astCtxSetter(
			doc,
			fieldSetter,
			fieldName,
			fieldType,
		)
		addDecl(astFile, setterDecl)
		ctx.Methods = append(ctx.Methods, setterDecl)
	}
}

// # gen segment: Singleton instance #
func genSingletonNewAst(ctx *model.Ctx, astFile *ast.File) {
	ctxVar := "ctx"

	var docs []*ast.Comment
	var stmts []ast.Stmt
	// [code] ctx := &Ctx{}
	stmts = append(stmts, astDefineStmt(
		ast.NewIdent(ctxVar),
		astDeclareRef(ast.NewIdent(CtxType), nil),
	))

	// create instances
	for _, instance := range ctx.SingletonInstances {
		instanceOrder := instance.GetOrder()
		instanceName := instance.GetInstance()
		instanceFunc := instance.GetFunc()
		instanceType := instance.GetType()

		fieldName := utils.FirstToLower(instanceName)
		fieldExpr := astSelectorExpr(ctxVar, fieldName)
		constructor := instance.GetConstructor()

		if constructor == nil {
			if instanceFunc.FuncName == "" {
				// [code] ctx.{{PrivateName}} = &{{Package}}.{{EventName}}{}
				constructor = astDeclareRef(
					instanceType,
					nil,
				)
			} else {
				// [code] ctx.{{PrivateName}} = {{Package}}.{{FuncName}}()
				constructor = &ast.CallExpr{Fun: astSelectorExpr(instanceFunc.Package, instanceFunc.FuncName)}
			}
		}

		stmts = append(stmts, astAssignStmt(
			fieldExpr,
			constructor,
		))

		if instanceOrder != "" {
			docs = append(docs, &ast.Comment{
				Text: "//  " + instanceOrder,
			})
		}
	}

	// call func
	for _, instance := range ctx.SingletonInstances {
		if _, ok := instance.(*model.Provide); ok {
			instanceName := instance.GetInstance()
			fieldName := utils.FirstToLower(instanceName)
			fieldExpr := astSelectorExpr(ctxVar, fieldName)

			handler := instance.GetHandler()
			if handler != "" {

				var instanceCallExpr *ast.CallExpr
				if strings.Contains(handler, ".") {
					// [code] {{Handler}}(ctx.{{}})
					instanceCallExpr = &ast.CallExpr{
						Fun:  ast.NewIdent(handler),
						Args: []ast.Expr{fieldExpr},
					}
				} else {
					// [code] ctx.{{Handler}}(ctx.{{}})
					instanceCallExpr = &ast.CallExpr{
						Fun:  astSelectorExpr(ctxVar, handler),
						Args: []ast.Expr{fieldExpr},
					}
				}

				// [code] {{Package}}.{{FunName}}(...)
				stmts = append(stmts, &ast.ExprStmt{
					X: instanceCallExpr,
				})
			}
		}
	}

	// [code] return ctx
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent(ctxVar),
		},
	})

	funcDecl := astFuncDecl(
		nil,
		"New",
		nil,
		[]*ast.Field{
			{
				Type: astStarExpr(ast.NewIdent(CtxType)),
			},
		},
		stmts,
	)
	if docs != nil {
		docs = append([]*ast.Comment{{Text: "// Action list:"}}, docs...)
		ctx.Doc = docs
		funcDecl.Doc = &ast.CommentGroup{List: docs}
	}

	addDecl(astFile, funcDecl)
}
