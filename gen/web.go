package gen

import (
	"github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
)

func genWebFile(moduleInfo *model.ModuleInfo, dir string) error {
	filename := filepath.Join(dir, global.GenWebFilename)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	if moduleInfo.WebAppInstances == nil {
		return nil
	}

	astFile := &ast.File{
		Name:  astIdent(global.GenPackage),
		Scope: ast.NewScope(nil),
	}

	genWebImportsAst(moduleInfo, astFile)

	genWebFuncAst(astFile)

	genWebAppStartupAst(moduleInfo, astFile)

	genMiddlewareAst(moduleInfo, astFile)

	genRouterAst(moduleInfo, astFile)

	err = format.Node(file, token.NewFileSet(), astFile)
	if err != nil {
		return err
	}
	return nil
}

func genWebImportsAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	astImport(astFile, "", "mime/multipart")
	astImport(astFile, "", "reflect")
	astImport(astFile, "", "strconv")
	astImport(astFile, "", "github.com/gofiber/fiber/v2")
	astImport(astFile, "", "fmt")
	for _, instance := range moduleInfo.WebAppInstances {
		astImport(astFile, "", instance.Import)
		for _, importInfo := range instance.Imports {
			importName := importInfo.Name
			if importName == "_" {
				importName = ""
			}
			astImport(astFile, importName, importInfo.Path)
		}

		for _, middleware := range instance.Middlewares {
			astImport(astFile, "", middleware.Import)
			for _, importInfo := range middleware.Imports {
				importName := importInfo.Name
				if importName == "_" {
					importName = ""
				}
				astImport(astFile, importName, importInfo.Path)
			}
		}

		for _, router := range instance.Routers {
			astImport(astFile, "", router.Import)
			for _, importInfo := range router.Imports {
				importName := importInfo.Name
				if importName == "_" {
					importName = ""
				}
				astImport(astFile, importName, importInfo.Path)
			}
		}
	}
	addImportDecl(astFile)
}

func genWebFuncAst(astFile *ast.File) {
	// [code] func Params(webCtx *fiber.Ctx, key string, defaultValue ...string) string
	addDecl(astFile, astFuncDecl(
		nil,
		"Params",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("string")}),
		},
		[]*ast.Field{
			astField("", astIdent("string")),
		},
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
	))

	// [code] func ParamsInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"ParamsInt",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
		},
		[]*ast.Field{
			astField("", astIdent("int")),
			astField("", astIdent("error")),
		},
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
	))

	// [code] func ParamsBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"ParamsBool",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("bool")}),
		},
		[]*ast.Field{
			astField("", astIdent("bool")),
			astField("", astIdent("error")),
		},
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
	))

	// [code] func ParamsFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"ParamsFloat",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("float64")}),
		},
		[]*ast.Field{
			astField("", astIdent("float64")),
			astField("", astIdent("error")),
		},
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
	))

	// [code] func ParamsParser(webCtx *fiber.Ctx, out any) error
	addDecl(astFile, astFuncDecl(
		nil,
		"ParamsParser",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("out", astIdent("any")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
		},
		[]*ast.Field{
			astField("", astIdent("error")),
		},
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
	))

	// [code] func Query(webCtx *fiber.Ctx, key string, defaultValue ...string) string
	addDecl(astFile, astFuncDecl(
		nil,
		"Query",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("string")}),
		},
		[]*ast.Field{
			astField("", astIdent("string")),
		},
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
	))
	// [code] func QueryInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"QueryInt",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
		},
		[]*ast.Field{
			astField("", astIdent("int")),
			astField("", astIdent("error")),
		},
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
	))
	// [code] func QueryBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"QueryBool",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("bool")}),
		},
		[]*ast.Field{
			astField("", astIdent("bool")),
			astField("", astIdent("error")),
		},
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
	))
	// [code] func QueryFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"QueryFloat",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("float64")}),
		},
		[]*ast.Field{
			astField("", astIdent("float64")),
			astField("", astIdent("error")),
		},
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
	))

	// [code] func QueryParser(webCtx *fiber.Ctx, out any) error
	addDecl(astFile, astFuncDecl(
		nil,
		"QueryParser",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("out", astIdent("any")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
		},
		[]*ast.Field{
			astField("", astIdent("error")),
		},
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
	))

	// [code] func Header(webCtx *fiber.Ctx, key string, defaultValue ...string) string
	addDecl(astFile, astFuncDecl(
		nil,
		"Header",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("string")}),
		},
		[]*ast.Field{
			astField("", astIdent("string")),
		},
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
	))
	// [code] func HeaderInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"HeaderInt",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
		},
		[]*ast.Field{
			astField("", astIdent("int")),
			astField("", astIdent("error")),
		},
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
	))
	// [code] func HeaderBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"HeaderBool",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("bool")}),
		},
		[]*ast.Field{
			astField("", astIdent("bool")),
			astField("", astIdent("error")),
		},
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
	))
	// [code] func HeaderFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"HeaderFloat",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("float64")}),
		},
		[]*ast.Field{
			astField("", astIdent("float64")),
			astField("", astIdent("error")),
		},
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
	))
	// [code] func HeaderParser(webCtx *fiber.Ctx, out any) error
	addDecl(astFile, astFuncDecl(
		nil,
		"HeaderParser",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("out", astIdent("any")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
		},
		[]*ast.Field{
			astField("", astIdent("error")),
		},
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
	))

	// [code] func FormString(webCtx *fiber.Ctx, key string, defaultValue ...string) string
	addDecl(astFile, astFuncDecl(
		nil,
		"FormString",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("string")}),
		},
		[]*ast.Field{
			astField("", astIdent("string")),
		},
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
	))
	// [code] func FormFile(webCtx *fiber.Ctx, key string) (*multipart.FileHeader, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"FormFile",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
		},
		[]*ast.Field{
			astField("", astStarExpr(astSelectorExpr("multipart", "FileHeader"))),
			astField("", astIdent("error")),
		},
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
	))
	// [code] func FormInt(webCtx *fiber.Ctx, key string, defaultValue ...int) (int, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"FormInt",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
		},
		[]*ast.Field{
			astField("", astIdent("int")),
			astField("", astIdent("error")),
		},
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
	))
	// [code] func FormBool(webCtx *fiber.Ctx, key string, defaultValue ...bool) (bool, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"FormBool",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("bool")}),
		},
		[]*ast.Field{
			astField("", astIdent("bool")),
			astField("", astIdent("error")),
		},
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
	))
	// [code] func FormFloat(webCtx *fiber.Ctx, key string, defaultValue ...float64) (float64, error)
	addDecl(astFile, astFuncDecl(
		nil,
		"FormFloat",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("key", astIdent("string")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("float64")}),
		},
		[]*ast.Field{
			astField("", astIdent("float64")),
			astField("", astIdent("error")),
		},
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
	))
	// [code] func FormParser(webCtx *fiber.Ctx, out any) error
	addDecl(astFile, astFuncDecl(
		nil,
		"FormParser",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("out", astIdent("any")),
			astField("defaultValue", &ast.Ellipsis{Elt: astIdent("int")}),
		},
		[]*ast.Field{
			astField("", astIdent("error")),
		},
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
	))

	// [code] func Body(webCtx *fiber.Ctx) []byte
	addDecl(astFile, astFuncDecl(
		nil,
		"Body",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
		},
		[]*ast.Field{
			astField("", &ast.ArrayType{
				Elt: astIdent("byte"),
			}),
		},
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExpr("webCtx", "Body"),
					},
				},
			},
		},
	))
	// [code] func BodyString(webCtx *fiber.Ctx) string
	addDecl(astFile, astFuncDecl(
		nil,
		"BodyString",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
		},
		[]*ast.Field{
			astField("", astIdent("string")),
		},
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
	))
	// [code] func BodyParser(webCtx *fiber.Ctx, out any) error
	addDecl(astFile, astFuncDecl(
		nil,
		"BodyParser",
		[]*ast.Field{
			astField("webCtx", astStarExpr(
				astSelectorExpr("fiber", "Ctx"))),
			astField("out", astIdent("any")),
		},
		[]*ast.Field{
			astField("", astIdent("error")),
		},
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
	))
}

func paramConvStmts(strCall *ast.CallExpr, convCall *ast.CallExpr) []ast.Stmt {
	return []ast.Stmt{
		astDefineStmt(astIdent("str"), strCall),
		// [code] if str == "" && defaultValue != nil
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				Op: token.LAND,
				X: &ast.BinaryExpr{
					Op: token.EQL,
					X:  astIdent("str"),
					Y:  astStringExpr(""),
				},
				Y: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  astIdent("defaultValue"),
					Y:  astIdent("nil"),
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					// [code] return defaultValue[0], nil
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.IndexExpr{
								Index: astIntExpr("0"),
								X:     astIdent("defaultValue"),
							},
							astIdent("nil"),
						},
					},
				},
			},
		},
		// [code] return strconv.ParseBool(str)
		&ast.ReturnStmt{
			Results: []ast.Expr{
				convCall,
			},
		},
	}
}

func errorReturnStmts() *ast.IfStmt {
	return &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			Op: token.NEQ,
			X:  astIdent("err"),
			Y:  astIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{astIdent("err")},
				},
			},
		},
	}
}
func genWebAppStartupAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	recvVar := utils.FirstToLower(global.StructName)

	for _, instance := range moduleInfo.WebAppInstances {
		params := make([]*ast.Field, 0)
		if instance.FuncName != "" {
			for _, paramInfo := range instance.Params {
				if paramInfo.Source == "" {
					// [code] {{ParamInstance}} {{ParamType}},
					paramInstance := paramInfo.Instance
					params = append(params,
						astField(paramInstance,
							utils.AccessType(
								paramInfo.Type,
								instance.Package,
								global.GenPackage,
							),
						),
					)
				}
			}

		} else {
			params = append(params,
				astField("host", astIdent("string")),
				astField("port", astIdent("uint")),
			)
		}
		stmts := make([]ast.Stmt, 0)

		for _, staticResource := range instance.Statics {
			// [code] ctx.{{WebApp}}.Static({{Path}}, {{Dirname}}, {{...}})
			args := []ast.Expr{
				astStringExpr(staticResource.Path),
				astStringExpr(staticResource.Dirname),
			}
			if len(staticResource.Features) > 0 ||
				staticResource.Index != "" ||
				staticResource.MaxAge != 0 {
				elts := make([]ast.Expr, 0)
				for _, feature := range staticResource.Features {
					elts = append(elts, &ast.KeyValueExpr{
						Key:   astIdent(feature),
						Value: astIdent("true"),
					})
				}
				if staticResource.Index != "" {
					elts = append(elts, &ast.KeyValueExpr{
						Key:   astIdent("Index"),
						Value: astStringExpr(staticResource.Index),
					})
				}
				if staticResource.MaxAge != 0 {
					elts = append(elts, &ast.KeyValueExpr{
						Key:   astIdent("MaxAge"),
						Value: astIntExpr(strconv.Itoa(staticResource.MaxAge)),
					})
				}
				args = append(args, astDeclareExpr(
					astSelectorExpr("fiber", "Static"),
					elts,
				))
			}

			stmts = append(stmts, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: astSelectorExprRecur(
						astSelectorExpr(recvVar,
							instance.WebApp,
						),
						"Static"),
					Args: args,
				},
			})
		}

		for _, middleware := range instance.Middlewares {
			// [code] ctx.{{WebApp}}.Group({{Path}}, {{Proxy}})
			stmts = append(stmts, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: astSelectorExprRecur(
						astSelectorExpr(recvVar, instance.WebApp),
						"Group",
					),
					Args: []ast.Expr{
						astStringExpr(middleware.Path),
						astSelectorExpr(recvVar, middleware.Proxy),
					},
				},
			})
		}

		for _, router := range instance.Routers {
			for _, method := range router.Methods {
				// [code] ctx.{{WebApp}}.{{Method}}({{Path}}, {{Proxy}})
				stmts = append(stmts, &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: astSelectorExprRecur(
							astSelectorExpr(recvVar, instance.WebApp),
							method,
						),
						Args: []ast.Expr{
							astStringExpr(router.Path),
							astSelectorExpr(recvVar, router.Proxy),
						},
					},
				})
			}
		}

		if instance.FuncName != "" {
			// [code] host, port, err := {{Package}}.{{FuncName}}(...)
			args := make([]ast.Expr, 0)
			for _, paramInfo := range instance.Params {
				paramInstance := paramInfo.Instance

				if paramInfo.Source == "inject" {
					if paramInstance == "Ctx" {
						// [code] ctx,
						args = append(args, astIdent(recvVar))
					} else {
						// [code] ctx.{{ParamInstance}},
						if !moduleInfo.HasInstance(paramInstance) {
							utils.Failuref("%s, \"%s\" No matching Instance, at %s()", paramInfo.Comment, paramInstance, instance.FuncName)
						}
						args = append(args, astSelectorExpr(recvVar, paramInstance))
					}
				} else {
					// [code] {{ParamInstance}},
					args = append(args, astIdent(paramInstance))
				}
			}

			stmts = append(stmts, astDefineStmtMany(
				[]ast.Expr{
					astIdent("host"),
					astIdent("port"),
					astIdent("err"),
				},
				&ast.CallExpr{
					Fun:  astSelectorExpr(instance.Package, instance.FuncName),
					Args: args,
				},
			))

			stmts = append(stmts, errorReturnStmts())
		}

		stmts = append(stmts, &ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.CallExpr{
					Fun: astSelectorExprRecur(
						astSelectorExpr(recvVar, instance.WebApp),
						"Listen",
					),
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: astSelectorExpr("fmt", "Sprintf"),
							Args: []ast.Expr{
								astStringExpr("%s:%d"),
								astIdent("host"),
								astIdent("port"),
							},
						},
					},
				},
			},
		})

		addDecl(astFile, astFuncDecl(
			[]*ast.Field{
				astField(recvVar, astStarExpr(astIdent(global.StructName))),
			},
			instance.Proxy,
			params,
			[]*ast.Field{
				astField("", astIdent("error")),
			},
			stmts,
		))
	}

}

func defineParamStmt(convFunc string, param *model.FieldInfo) ast.Stmt {
	return astDefineStmt(
		astIdent(param.Instance),
		&ast.CallExpr{
			Fun: astIdent(convFunc),
			Args: []ast.Expr{
				astIdent("webCtx"),
				astStringExpr(param.Instance),
			},
		},
	)
}

func defineParamWithError(convFunc string, param *model.FieldInfo) []ast.Stmt {
	return []ast.Stmt{
		astDefineStmtMany(
			[]ast.Expr{
				astIdent(param.Instance),
				astIdent("err"),
			},
			&ast.CallExpr{
				Fun: astIdent(convFunc),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astStringExpr(param.Instance),
				},
			},
		),
		errorReturnStmts(),
	}
}

func defineParamByParser(convFunc string, param *model.FieldInfo, packageInfo *model.PackageInfo) []ast.Stmt {
	return []ast.Stmt{
		// [code] {{ParamInstance}} := &{{Package}}.{{ParamType}}{}
		astDefineStmt(
			astIdent(param.Instance),
			astDeclareRef(
				utils.AccessType(
					param.Type,
					packageInfo.Package,
					global.GenPackage,
				),
				nil,
			),
		),
		// [code] err := {{ParamSource}}Parser(webCtx, {{ParamInstance}})
		astAssignStmt(
			astIdent("err"),
			&ast.CallExpr{
				Fun: astIdent(convFunc),
				Args: []ast.Expr{
					astIdent("webCtx"),
					astIdent(param.Instance),
				},
			},
		),
		errorReturnStmts(),
	}
}
func genWebBodyParam(bodyParam *model.FieldInfo, packageInfo *model.PackageInfo, funcName string) []ast.Stmt {
	if bodyParam != nil {
		bodyAstType := bodyParam.Type
		switch bodyAstType.(type) {
		case *ast.ArrayType:
			byteArr := bodyAstType.(*ast.ArrayType)
			if byteArr.Elt.(*ast.Ident).String() == "byte" {
				return []ast.Stmt{
					astDefineStmt(
						astIdent(bodyParam.Instance),
						&ast.CallExpr{
							Fun: astIdent("Body"),
							Args: []ast.Expr{
								astIdent("webCtx"),
							},
						},
					),
				}
			}
			break
		case *ast.Ident:
			typeIdent := bodyAstType.(*ast.Ident).String()
			if typeIdent == "string" {
				return []ast.Stmt{
					astDefineStmt(
						astIdent(bodyParam.Instance),
						&ast.CallExpr{
							Fun: astIdent("BodyString"),
							Args: []ast.Expr{
								astIdent("webCtx"),
							},
						},
					),
				}
			} else if utils.IsFirstLower(typeIdent) {
				utils.Failuref("%s, unsupport type %s, at %s()", bodyParam.Comment, utils.TypeToString(bodyAstType), funcName)
			}
			break
		}
		return defineParamByParser("BodyParser", bodyParam, packageInfo)
	}
	return nil
}

func genWebHeaderParams(headerParams []*model.FieldInfo, packageInfo *model.PackageInfo, funcName string) []ast.Stmt {
	if len(headerParams) > 0 {
		stmts := make([]ast.Stmt, 0)
		for _, param := range headerParams {
			paramType := param.Type
			switch paramType.(type) {
			case *ast.Ident:
				typeStr := paramType.(*ast.Ident).String()
				switch typeStr {
				case "string":
					stmts = append(stmts, defineParamStmt("Header", param))
					continue
				case "int":
					stmts = append(stmts, defineParamWithError("HeaderInt", param)...)
					continue
				case "bool":
					stmts = append(stmts, defineParamWithError("HeaderBool", param)...)
					continue
				case "float64":
					stmts = append(stmts, defineParamWithError("HeaderFloat", param)...)
					continue
				default:
					if utils.IsFirstLower(typeStr) {
						utils.Failuref("%s, unsupport type %s, at %s()", param.Comment, utils.TypeToString(paramType), funcName)
					}
				}
				break
			}
			stmts = append(stmts, defineParamByParser("HeaderParser", param, packageInfo)...)
		}
		return stmts
	}
	return nil
}
func genWebQueryParams(queryParams []*model.FieldInfo, packageInfo *model.PackageInfo, funcName string) []ast.Stmt {
	if len(queryParams) > 0 {
		stmts := make([]ast.Stmt, 0)
		for _, param := range queryParams {
			paramType := param.Type
			switch paramType.(type) {
			case *ast.Ident:
				typeStr := paramType.(*ast.Ident).String()
				switch typeStr {
				case "string":
					stmts = append(stmts, defineParamStmt("Query", param))
					continue
				case "int":
					stmts = append(stmts, defineParamWithError("QueryInt", param)...)
					continue
				case "bool":
					stmts = append(stmts, defineParamWithError("QueryBool", param)...)
					continue
				case "float64":
					stmts = append(stmts, defineParamWithError("QueryFloat", param)...)
					continue
				default:
					if utils.IsFirstLower(typeStr) {
						utils.Failuref("%s, unsupport type %s, at %s()", param.Comment, utils.TypeToString(paramType), funcName)
					}
				}
				break
			}
			stmts = append(stmts, defineParamByParser("QueryParser", param, packageInfo)...)
		}
		return stmts
	}
	return nil
}

func genWebPathParams(pathParams []*model.FieldInfo, packageInfo *model.PackageInfo, funcName string) []ast.Stmt {
	if len(pathParams) > 0 {
		stmts := make([]ast.Stmt, 0)
		for _, param := range pathParams {
			paramType := param.Type
			switch paramType.(type) {
			case *ast.Ident:
				typeStr := paramType.(*ast.Ident).String()
				switch typeStr {
				case "string":
					stmts = append(stmts, defineParamStmt("Params", param))
					continue
				case "int":
					stmts = append(stmts, defineParamWithError("ParamsInt", param)...)
					continue
				case "bool":
					stmts = append(stmts, defineParamWithError("ParamsBool", param)...)
					continue
				case "float64":
					stmts = append(stmts, defineParamWithError("ParamsFloat", param)...)
					continue
				default:
					if utils.IsFirstLower(typeStr) {
						utils.Failuref("%s, unsupport type %s, at %s()", param.Comment, utils.TypeToString(paramType), funcName)
					}
				}
				break
			}
			stmts = append(stmts, defineParamByParser("ParamsParser", param, packageInfo)...)
		}
		return stmts
	}
	return nil
}
func genWebFormParams(formParams []*model.FieldInfo, packageInfo *model.PackageInfo, funcName string) []ast.Stmt {
	if len(formParams) > 0 {
		stmts := make([]ast.Stmt, 0)
		for _, param := range formParams {
			paramType := param.Type
			switch paramType.(type) {
			case *ast.Ident:
				typeStr := paramType.(*ast.Ident).String()
				switch typeStr {
				case "string":
					stmts = append(stmts, defineParamStmt("FormString", param))
					continue
				case "int":
					stmts = append(stmts, defineParamWithError("FormInt", param)...)
					continue
				case "bool":
					stmts = append(stmts, defineParamWithError("FormBool", param)...)
					continue
				case "float64":
					stmts = append(stmts, defineParamWithError("FormFloat", param)...)
					continue
				default:
					if utils.IsFirstLower(typeStr) {
						utils.Failuref("%s, unsupport type %s, at %s()", param.Comment, utils.TypeToString(paramType), funcName)
					}
				}
				break
			case *ast.StarExpr:
				// [code] *multipart.FileHeader
				starX := paramType.(*ast.StarExpr).X
				selectorExpr, ok := starX.(*ast.SelectorExpr)
				if ok {
					selectorX, ok := selectorExpr.X.(*ast.Ident)
					if ok {
						if selectorX.String() == "multipart" && selectorExpr.Sel.String() == "FileHeader" {
							stmts = append(stmts, defineParamStmt("FormFile", param))
							continue
						}
					}
				}
				break
			}
			stmts = append(stmts, defineParamByParser("FormParser", param, packageInfo)...)
		}
		return stmts
	}
	return nil
}
