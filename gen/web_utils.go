package gen

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
)

func genWebUtilsFile(moduleInfo *model.ModuleInfo, dir string) error {
	fileDir := filepath.Join(dir, GenUtilsPackage)
	filename := filepath.Join(fileDir, GenWebUtilsFilename)

	if moduleInfo.WebAppInstances == nil {
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
		Name:  astIdent(GenUtilsPackage),
		Scope: ast.NewScope(nil),
	}

	genWebUtilsImportsAst(moduleInfo, astFile)

	genWebUtilsFuncAst(moduleInfo, astFile)

	astFile, err := utils.OptimizeCode(filename, astFile, moduleInfo,
		"// Code generated by \"inject-golang -m web\"; DO NOT EDIT.")
	if err != nil {
		return err
	}

	return utils.GenerateCode(filename, astFile, moduleInfo)
}

func genWebUtilsImportsAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	addImport(astFile, moduleInfo, "", "mime/multipart")
	addImport(astFile, moduleInfo, "", "reflect")
	addImport(astFile, moduleInfo, "", "strconv")
	addImport(astFile, moduleInfo, "", "github.com/gofiber/fiber/v2")
	addImport(astFile, moduleInfo, "", "fmt")
}

func genWebUtilsFuncAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	funcDoc := &ast.Comment{
		Text: "// Generate by system",
	}

	// [code] func Params(webCtx *fiber.Ctx, key string, defaultValue ...string) string
	funcName := "Params"
	funcParams := []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("string")}),
	}
	funcResults := []*ast.Field{
		astField("", astIdent("string")),
	}
	funcDecl := astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "Params"),
						Args: []ast.Expr{
							astIdent("key"),
							astIdent("defaultValue"),
						},
						Ellipsis: 1,
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func ParamsInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error)
	funcName = "ParamsInt"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("int")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "ParamsInt"),
						Args: []ast.Expr{
							astIdent("key"),
							astIdent("defaultValue"),
						},
						Ellipsis: 1,
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func ParamsBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error)
	funcName = "ParamsBool"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("bool")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("bool")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("Params"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "ParseBool"),
				Args: []ast.Expr{
					astIdent("str"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func ParamsFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error)
	funcName = "ParamsFloat"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("float64")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("float64")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("Params"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "ParseFloat"),
				Args: []ast.Expr{
					astIdent("str"),
					astIntExpr("64"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func ParamsParser(webCtx *fiber.Ctx, out any) error
	funcName = "ParamsParser"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("out", astIdent("any")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "ParamsParser"),
						Args: []ast.Expr{
							astIdent("out"),
						},
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func Query(webCtx *fiber.Ctx, key string, defaultValue ...string) string
	funcName = "Query"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("string")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("string")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "Query"),
						Args: []ast.Expr{
							astIdent("key"),
							astIdent("defaultValue"),
						},
						Ellipsis: 1,
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func QueryInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error)
	funcName = "QueryInt"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("int")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("Query"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "Atoi"),
				Args: []ast.Expr{
					astIdent("str"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func QueryBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error)
	funcName = "QueryBool"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("bool")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("bool")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("Query"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "ParseBool"),
				Args: []ast.Expr{
					astIdent("str"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func QueryFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error)
	funcName = "QueryFloat"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("float64")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("float64")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("Query"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "ParseFloat"),
				Args: []ast.Expr{
					astIdent("str"),
					astIntExpr("64"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func QueryParser(webCtx *fiber.Ctx, out any) error
	funcName = "QueryParser"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("out", astIdent("any")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "QueryParser"),
						Args: []ast.Expr{
							astIdent("out"),
						},
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func Header(webCtx *fiber.Ctx, key string, defaultValue ...string) string
	funcName = "Header"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("string")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("string")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "GetRespHeader"),
						Args: []ast.Expr{
							astIdent("key"),
							astIdent("defaultValue"),
						},
						Ellipsis: 1,
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func HeaderInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error)
	funcName = "HeaderInt"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("int")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("Header"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "Atoi"),
				Args: []ast.Expr{
					astIdent("str"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func HeaderBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error)
	funcName = "HeaderBool"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("bool")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("bool")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("Header"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "ParseBool"),
				Args: []ast.Expr{
					astIdent("str"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func HeaderFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error)
	funcName = "HeaderFloat"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("float64")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("float64")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("Header"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "ParseFloat"),
				Args: []ast.Expr{
					astIdent("str"),
					astIntExpr("64"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func HeaderParser(webCtx *fiber.Ctx, out any) error

	funcName = "HeaderParser"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("out", astIdent("any")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "ReqHeaderParser"),
						Args: []ast.Expr{
							astIdent("out"),
						},
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func FormString(webCtx *fiber.Ctx, key string, defaultValue ...string) string
	funcName = "FormString"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("string")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("string")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "FormValue"),
						Args: []ast.Expr{
							astIdent("key"),
							astIdent("defaultValue"),
						},
						Ellipsis: 1,
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func FormFile(webCtx *fiber.Ctx, key string) (*multipart.FileHeader, error)
	funcName = "FormFile"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
	}
	funcResults = []*ast.Field{
		astField("", astStarExpr(astSelectorExpr("multipart", "FileHeader"))),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "FormFile"),
						Args: []ast.Expr{
							astIdent("key"),
						},
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func FormInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error)
	funcName = "FormInt"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("int")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("FormString"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "Atoi"),
				Args: []ast.Expr{
					astIdent("str"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func FormBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error)
	funcName = "FormBool"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("bool")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("bool")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("FormString"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "ParseBool"),
				Args: []ast.Expr{
					astIdent("str"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func FormFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error)
	funcName = "FormFloat"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("key", astIdent("string")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("float64")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("float64")),
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		paramConvStmts(
			&ast.CallExpr{
				Fun: astIdent("FormString"),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent("key"),
				},
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("strconv", "ParseFloat"),
				Args: []ast.Expr{
					astIdent("str"),
					astIntExpr("64"),
				},
			},
		),
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func FormParser(webCtx *fiber.Ctx, out any) error
	funcName = "FormParser"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("out", astIdent("any")),
		astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			astDefineStmt(
				astIdent("elem"),
				&ast.CallExpr{
					Fun: astSelectorExprRecur(
						&ast.CallExpr{
							Fun: astSelectorExpr(
								"reflect",
								"ValueOf",
							),
							Args: []ast.Expr{
								astIdent("out"),
							},
						},
						"Elem",
					),
				},
			),
			astDefineStmtMany(
				[]ast.Expr{
					astIdent("form"),
					astIdent("err"),
				},
				&ast.CallExpr{
					Fun: astSelectorExpr(
						"webCtx",
						"MultipartForm",
					),
				},
			),
			errorReturnStmts(),
			&ast.RangeStmt{
				Key:   astIdent("key"),
				Value: astIdent("strArr"),
				Tok:   token.DEFINE,
				X:     astSelectorExpr("form", "Value"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.RangeStmt{
							Key:   astIdent("_"),
							Value: astIdent("value"),
							Tok:   token.DEFINE,
							X:     astIdent("strArr"),
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									astDefineStmt(
										astIdent("field"),
										&ast.CallExpr{
											Fun: astSelectorExpr("elem", "FieldByName"),
											Args: []ast.Expr{
												astIdent("key"),
											},
										},
									),
									&ast.IfStmt{
										Cond: &ast.BinaryExpr{
											Op: token.LAND,
											X: &ast.CallExpr{
												Fun: astSelectorExpr("field", "IsValid"),
											},
											Y: &ast.CallExpr{
												Fun: astSelectorExpr("field", "CanSet"),
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.SwitchStmt{
													Tag: &ast.CallExpr{
														Fun: astSelectorExpr("field", "Kind"),
													},
													Body: &ast.BlockStmt{
														List: []ast.Stmt{
															&ast.CaseClause{
																List: []ast.Expr{
																	astSelectorExpr("reflect", "String"),
																},
																Body: []ast.Stmt{
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: astSelectorExpr("field", "SetString"),
																			Args: []ast.Expr{
																				astIdent("value"),
																			},
																		},
																	},
																	&ast.BranchStmt{Tok: token.BREAK},
																},
															},
															&ast.CaseClause{
																List: []ast.Expr{
																	astSelectorExpr("reflect", "Int"),
																},
																Body: []ast.Stmt{
																	astDefineStmtMany(
																		[]ast.Expr{
																			astIdent("intValue"),
																			astIdent("err"),
																		},
																		&ast.CallExpr{
																			Fun: astSelectorExpr("strconv", "Atoi"),
																			Args: []ast.Expr{
																				astIdent("value"),
																			},
																		},
																	),
																	errorReturnStmts(),
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: astSelectorExpr("field", "SetInt"),
																			Args: []ast.Expr{
																				&ast.CallExpr{
																					Fun: astIdent("int64"),
																					Args: []ast.Expr{
																						astIdent("intValue"),
																					},
																				},
																			},
																		},
																	},
																	&ast.BranchStmt{Tok: token.BREAK},
																},
															},
															&ast.CaseClause{
																List: []ast.Expr{
																	astSelectorExpr("reflect", "Bool"),
																},
																Body: []ast.Stmt{
																	astDefineStmtMany(
																		[]ast.Expr{
																			astIdent("boolValue"),
																			astIdent("err"),
																		},
																		&ast.CallExpr{
																			Fun: astSelectorExpr("strconv", "ParseBool"),
																			Args: []ast.Expr{
																				astIdent("value"),
																			},
																		},
																	),
																	errorReturnStmts(),
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: astSelectorExpr("field", "SetBool"),
																			Args: []ast.Expr{
																				astIdent("boolValue"),
																			},
																		},
																	},
																	&ast.BranchStmt{Tok: token.BREAK},
																},
															},
															&ast.CaseClause{
																List: []ast.Expr{
																	astSelectorExpr("reflect", "Float64"),
																},
																Body: []ast.Stmt{
																	astDefineStmtMany(
																		[]ast.Expr{
																			astIdent("floatValue"),
																			astIdent("err"),
																		},
																		&ast.CallExpr{
																			Fun: astSelectorExpr("strconv", "ParseFloat"),
																			Args: []ast.Expr{
																				astIdent("value"),
																				astIntExpr("64"),
																			},
																		},
																	),
																	errorReturnStmts(),
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: astSelectorExpr("field", "SetFloat"),
																			Args: []ast.Expr{
																				astIdent("floatValue"),
																			},
																		},
																	},
																	&ast.BranchStmt{Tok: token.BREAK},
																},
															},
															&ast.CaseClause{
																Body: []ast.Stmt{
																	&ast.ReturnStmt{
																		Results: []ast.Expr{
																			&ast.CallExpr{
																				Fun: astSelectorExpr("fmt", "Errorf"),
																				Args: []ast.Expr{
																					astStringExpr("unsupported type %T"),
																					astIdent("value"),
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
										Else: &ast.IfStmt{
											Cond: &ast.UnaryExpr{
												Op: token.NOT,
												X: &ast.CallExpr{
													Fun: astSelectorExpr("field", "IsValid"),
												},
											},
											Body: &ast.BlockStmt{},
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.RangeStmt{
				Key:   astIdent("key"),
				Value: astIdent("fileArr"),
				Tok:   token.DEFINE,
				X:     astSelectorExpr("form", "File"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.RangeStmt{
							Key:   astIdent("_"),
							Value: astIdent("file"),
							Tok:   token.DEFINE,
							X:     astIdent("fileArr"),
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									astDefineStmt(
										astIdent("field"),
										&ast.CallExpr{
											Fun: astSelectorExpr("elem", "FieldByName"),
											Args: []ast.Expr{
												astIdent("key"),
											},
										},
									),
									&ast.IfStmt{
										Cond: &ast.BinaryExpr{
											Op: token.LAND,
											X: &ast.CallExpr{
												Fun: astSelectorExpr("field", "IsValid"),
											},
											Y: &ast.CallExpr{
												Fun: astSelectorExpr("field", "CanSet"),
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: astSelectorExpr("field", "Set"),
														Args: []ast.Expr{
															&ast.CallExpr{
																Fun: astSelectorExpr("reflect", "ValueOf"),
																Args: []ast.Expr{
																	astIdent("file"),
																},
															},
														},
													},
												},
											},
										},
										Else: &ast.IfStmt{
											Cond: &ast.UnaryExpr{
												Op: token.NOT,
												X: &ast.CallExpr{
													Fun: astSelectorExpr("field", "IsValid"),
												},
											},
											Body: &ast.BlockStmt{},
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.ReturnStmt{
				Results: []ast.Expr{
					astIdent("nil"),
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func Body(webCtx *fiber.Ctx) []byte
	funcName = "Body"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
	}
	funcResults = []*ast.Field{
		astField("", &ast.ArrayType{
			Elt: astIdent("byte"),
		}),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "Body"),
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func BodyString(webCtx *fiber.Ctx) string
	funcName = "BodyString"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("string")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astIdent("string"),
						Args: []ast.Expr{
							&ast.CallExpr{
								Fun: astIdent("Body"),
								Args: []ast.Expr{
									astIdent("webCtx"),
								},
							},
						},
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)

	// [code] func BodyParser(webCtx *fiber.Ctx, out any) error
	funcName = "BodyParser"
	funcParams = []*ast.Field{
		astField("webCtx", astStarExpr(
			astSelectorExpr("fiber", "Ctx"))),
		astField("out", astIdent("any")),
	}
	funcResults = []*ast.Field{
		astField("", astIdent("error")),
	}
	funcDecl = astFuncDecl(
		nil,
		funcName,
		funcParams,
		funcResults,
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "BodyParser"),
						Args: []ast.Expr{
							astIdent("out"),
						},
					},
				},
			},
		},
	)
	funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: fmt.Sprintf("// %s", funcName),
		},
		funcDoc,
	}}
	addDecl(astFile, funcDecl)
}
