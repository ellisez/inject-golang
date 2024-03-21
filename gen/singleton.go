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

	if ctx.SingletonInstance.Len() == 0 {
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

	genSingletonImportsAst(ctx, astFile, filename)

	genSingletonStructAst(ctx, astFile)

	genSingletonGetterAndSetterAst(ctx, astFile)

	genCtxNewAst(ctx, astFile)

	genSingletonNewAst(ctx, astFile)

	return utils.GenerateCode(filename, astFile, ctx,
		`// Code generated by "inject-golang -m singleton"; DO NOT EDIT.`)
}

func genSingletonImportsAst(ctx *model.Ctx, astFile *ast.File, filename string) {

	for i := 0; i < ctx.SingletonInstance.Len(); i++ {
		instance, webApplication := ctx.SingletonInstance.IndexOf(i)
		if webApplication == nil {
			for _, importNode := range instance.Imports {
				importName := importNode.Alias
				if importName == "_" {
					importName = ""
				}
				err := addImport(astFile, ctx, importName, importNode.Path)
				if err != nil {
					utils.Failuref("%s, %s", filename, err.Error())
				}
			}
		}
	}
	if ctx.SingletonInstance.WebLen() > 0 {
		err := addImport(astFile, ctx, "", "github.com/gofiber/fiber/v2")
		if err != nil {
			utils.Failuref("%s, %s", filename, err.Error())
		}
	}
}

// # gen segment: Struct #
func genSingletonStructAst(ctx *model.Ctx, astFile *ast.File) {
	var fields []*ast.Field
	for i := 0; i < ctx.SingletonInstance.Len(); i++ {
		instance, _ := ctx.SingletonInstance.IndexOf(i)
		instanceName := instance.Instance
		instanceType := instance.Type

		fieldName := utils.FirstToLower(instanceName)

		fieldType := instanceType

		field := astField(
			fieldName,
			fieldType,
		)
		fields = append(fields, field)
	}

	fields = append(fields, astField(ArgumentVar, &ast.MapType{
		Key:   ast.NewIdent("string"),
		Value: ast.NewIdent("any"),
	}))

	structDecl := astStructDecl(
		CtxType,
		fields,
	)

	addDecl(astFile, structDecl)
}

func genSingletonGetterAndSetterAst(ctx *model.Ctx, astFile *ast.File) {
	/// Getter / Setter
	for i := 0; i < ctx.SingletonInstance.Len(); i++ {
		instance, _ := ctx.SingletonInstance.IndexOf(i)
		instanceName := instance.Instance
		instanceType := instance.Type
		instanceFunc := instance.Func

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
		ctx.Methods[fieldGetter] = getterDecl

		setterDecl := astCtxSetter(
			doc,
			fieldSetter,
			fieldName,
			fieldType,
		)
		addDecl(astFile, setterDecl)
		ctx.Methods[fieldSetter] = setterDecl
	}
}

// # gen segment: SingletonInstance instance #
func genCtxNewAst(ctx *model.Ctx, astFile *ast.File) {
	ctxVar := "ctx"

	var docs []*ast.Comment
	var stmts []ast.Stmt
	// [code] ctx := &Ctx{}
	stmts = append(stmts, astDefineStmt(
		ast.NewIdent(ctxVar),
		astDeclareRef(ast.NewIdent(CtxType), nil),
	))
	// create args
	if ctx.SingletonInstance.ArgumentLen() > 0 {
		// [code] ctx.__args = map[string]any{}
		stmts = append(stmts, astAssignStmt(
			astSelectorExpr(ctxVar, ArgumentVar),
			astDeclareExpr(&ast.MapType{
				Key:   ast.NewIdent("string"),
				Value: ast.NewIdent("any"),
			}, nil),
		))
	}

	// create instances
	for i := 0; i < ctx.SingletonInstance.Len(); i++ {
		instance, _ := ctx.SingletonInstance.IndexOf(i)
		instanceOrder := instance.Order
		instanceName := instance.Instance
		instanceFunc := instance.Func
		instanceType := instance.Type

		fieldName := utils.FirstToLower(instanceName)
		var fieldExpr ast.Expr = astSelectorExpr(ctxVar, fieldName)
		constructor := instance.Constructor

		if instance.Mode == "argument" {
			fieldExpr = &ast.IndexExpr{
				X:     astSelectorExpr(ctxVar, ArgumentVar),
				Index: astStringExpr(instanceName),
			}
		}
		if constructor == nil {
			if instanceFunc.FuncName == "" {
				// [code] ctx.{{PrivateName}} = &{{Package}}.{{TypeName}}{}
				constructor = astDeclareRef(
					instanceType,
					nil,
				)
			} else {
				// [code] ctx.{{PrivateName}} = ctx.New{{Instance}}()
				constructor = &ast.CallExpr{Fun: astSelectorExpr(ctxVar, "New"+instanceName)}
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
	for i := 0; i < ctx.SingletonInstance.Len(); i++ {
		instance, webApplication := ctx.SingletonInstance.IndexOf(i)
		if webApplication == nil {

			handler := instance.Handler
			if handler != "" {

				var instanceCallExpr *ast.CallExpr
				if strings.Contains(handler, ".") {
					// [code] {{Handler}}()
					instanceCallExpr = &ast.CallExpr{
						Fun: ast.NewIdent(handler),
					}
				} else {
					// [code] ctx.{{Handler}}()
					instanceCallExpr = &ast.CallExpr{
						Fun: astSelectorExpr(ctxVar, handler),
					}
				}

				// [code] {{Package}}.{{FunName}}(...)
				stmts = append(stmts, &ast.ExprStmt{
					X: instanceCallExpr,
				})
			}
		}
	}

	if ctx.SingletonInstance.ArgumentLen() > 0 {
		stmts = append(stmts, astAssignStmt(
			astSelectorExpr(ctxVar, ArgumentVar),
			ast.NewIdent("nil"),
		))
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

func genSingletonNewAst(ctx *model.Ctx, astFile *ast.File) {
	ctxVar := "ctx"

	for i := 0; i < ctx.SingletonInstance.Len(); i++ {
		instance, webApplication := ctx.SingletonInstance.IndexOf(i)
		if webApplication == nil {
			instanceName := instance.Instance
			instanceFunc := instance.Func

			var stmts []ast.Stmt
			instanceCallExpr, varDefineStmts := astInstanceCallExpr(astSelectorExpr(instanceFunc.Package, instanceFunc.FuncName), instanceFunc, ctx, ctxVar)
			if varDefineStmts != nil {
				stmts = append(stmts, varDefineStmts...)
			}
			stmts = append(stmts, &ast.ReturnStmt{
				Results: []ast.Expr{
					instanceCallExpr,
				},
			})

			// [code] func (ctx *Container) {{Proxy}}(
			funcDecl := astInstanceProxyFunc(instanceFunc, "New"+instanceName)
			funcDecl.Body = &ast.BlockStmt{
				List: stmts,
			}
			funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{{
				Text: fmt.Sprintf("// Generate by annotations from %s.%s", instanceFunc.Package, instanceFunc.FuncName),
			}}}
			addDecl(astFile, funcDecl)
			ctx.Methods[funcDecl.Name.String()] = funcDecl
		}
	}
}
