package gen

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func genWebFile(ctx *model.Ctx, dir string) error {
	fileDir := filepath.Join(dir, GenInternalPackage)
	filename := filepath.Join(fileDir, GenWebFilename)

	if !ctx.HasWebInstance {
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

	genWebImportsAst(ctx, astFile)

	genWebAppStartupAst(ctx, astFile)

	genMiddlewareAst(ctx, astFile)

	genRouterAst(ctx, astFile)

	return utils.GenerateCode(filename, astFile, ctx,
		`// Code generated by "inject-golang -m web"; DO NOT EDIT.`)
}

func genWebImportsAst(ctx *model.Ctx, astFile *ast.File) {
	addImport(astFile, ctx, "", "fmt")
	addImport(astFile, ctx, "", "github.com/gofiber/fiber/v2")
	addImport(astFile, ctx, "", path.Join(Mod.Package, GenPackage, "utils"))

	for _, instance := range ctx.SingletonInstances {
		if webInstance, ok := instance.(*model.WebInstance); ok {
			for _, importInfo := range webInstance.Imports {
				importName := importInfo.Name
				if importName == "_" {
					importName = ""
				}
				addImport(astFile, ctx, importName, importInfo.Path)
			}

			for _, middleware := range webInstance.Middlewares {
				for _, importInfo := range middleware.Imports {
					importName := importInfo.Name
					if importName == "_" {
						importName = ""
					}
					addImport(astFile, ctx, importName, importInfo.Path)
				}
			}

			for _, router := range webInstance.Routers {
				for _, importInfo := range router.Imports {
					importName := importInfo.Name
					if importName == "_" {
						importName = ""
					}
					addImport(astFile, ctx, importName, importInfo.Path)
				}
			}
		}

	}
}

func errorReturnStmts() *ast.IfStmt {
	return &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			Op: token.NEQ,
			X:  ast.NewIdent("err"),
			Y:  ast.NewIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{ast.NewIdent("err")},
				},
			},
		},
	}
}
func genWebAppStartupAst(ctx *model.Ctx, astFile *ast.File) {
	ctxVar := utils.FirstToLower(CtxType)

	for _, instance := range ctx.SingletonInstances {
		if webInstance, ok := instance.(*model.WebInstance); ok {
			instanceName := webInstance.Instance
			instanceFunc := webInstance.Func
			instanceVar := utils.FirstToLower(CtxType)

			var proxyParams []*ast.Field
			if instanceFunc.FuncName == "" {
				// [code] host string, port uint,
				proxyParams = []*ast.Field{
					astField("host", ast.NewIdent("string")),
					astField("port", ast.NewIdent("uint")),
				}
			} else {
				// [code] {{ParamInstance}} {{ParamType}},
				proxyParams = astInstanceProxyParams(instance.GetFunc())
			}

			var stmts []ast.Stmt

			for _, resource := range webInstance.Resources {
				// [code] ctx.{{instance}}().Static({{Path}}, {{Path}}, {{...}})
				args := []ast.Expr{
					astStringExpr(resource.Path),
					astStringExpr(resource.Dirname),
				}
				if len(resource.Features) > 0 ||
					resource.Index != "" ||
					resource.MaxAge != 0 {
					var eltExpr []ast.Expr
					for _, feature := range resource.Features {
						eltExpr = append(eltExpr, &ast.KeyValueExpr{
							Key:   ast.NewIdent(feature),
							Value: ast.NewIdent("true"),
						})
					}
					if resource.Index != "" {
						eltExpr = append(eltExpr, &ast.KeyValueExpr{
							Key:   ast.NewIdent("Index"),
							Value: astStringExpr(resource.Index),
						})
					}
					if resource.MaxAge != 0 {
						eltExpr = append(eltExpr, &ast.KeyValueExpr{
							Key:   ast.NewIdent("MaxAge"),
							Value: astIntExpr(strconv.Itoa(resource.MaxAge)),
						})
					}
					args = append(args, astDeclareExpr(
						astSelectorExpr("fiber", "Static"),
						eltExpr,
					))
				}

				stmts = append(stmts, &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: astSelectorExprRecur(
							&ast.CallExpr{
								Fun: astSelectorExpr(ctxVar, webInstance.Instance),
							},
							"Static"),
						Args: args,
					},
				})
			}

			for _, middleware := range webInstance.Middlewares {
				// [code] ctx.{{instance}}().Group({{Path}}, {{Proxy}})
				stmts = append(stmts, &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: astSelectorExprRecur(
							&ast.CallExpr{
								Fun: astSelectorExpr(ctxVar, webInstance.Instance),
							},
							"Group",
						),
						Args: []ast.Expr{
							astStringExpr(middleware.Path),
							astSelectorExpr(ctxVar, middleware.FuncName),
						},
					},
				})
			}

			for _, router := range webInstance.Routers {
				for _, method := range router.Methods {
					// [code] ctx.{{instance}}().{{Method}}({{Path}}, {{Proxy}})
					stmts = append(stmts, &ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: astSelectorExprRecur(
								&ast.CallExpr{
									Fun: astSelectorExpr(ctxVar, webInstance.Instance),
								},
								method,
							),
							Args: []ast.Expr{
								astStringExpr(router.Path),
								astSelectorExpr(ctxVar, router.FuncName),
							},
						},
					})
				}
			}

			if webInstance.FuncName != "" {
				// [code] host, port, err := {{Package}}.{{FuncName}}(...)
				instanceCallExpr := astInstanceCallExpr(astSelectorExpr(webInstance.Package, webInstance.FuncName), webInstance.Func, ctx, ctxVar)

				stmts = append(stmts, astDefineStmtMany(
					[]ast.Expr{
						ast.NewIdent("host"),
						ast.NewIdent("port"),
						ast.NewIdent("err"),
					},
					instanceCallExpr,
				))

				stmts = append(stmts, errorReturnStmts())
			}

			stmts = append(stmts, &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: astSelectorExprRecur(
							&ast.CallExpr{
								Fun: astSelectorExpr(ctxVar, instanceName),
							},
							"Listen",
						),
						Args: []ast.Expr{
							&ast.CallExpr{
								Fun: astSelectorExpr("fmt", "Sprintf"),
								Args: []ast.Expr{
									astStringExpr("%s:%d"),
									ast.NewIdent("host"),
									ast.NewIdent("port"),
								},
							},
						},
					},
				},
			})

			doc := "// Generate by system"
			if instanceFunc.FuncName != "" {
				doc = fmt.Sprintf("// Generate by annotations from %s.%s", instanceFunc.Package, instanceFunc.FuncName)
			}

			proxyFuncDecl := astFuncDecl(
				[]*ast.Field{
					astField(instanceVar, astStarExpr(ast.NewIdent(CtxType))),
				},
				instanceName+"Startup",
				proxyParams,
				[]*ast.Field{
					astField("", ast.NewIdent("error")),
				},
				stmts,
			)
			proxyFuncDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{{
				Text: doc,
			}}}
			addDecl(astFile, proxyFuncDecl)
			ctx.Methods = append(ctx.Methods, proxyFuncDecl)
		}
	}

}

func defineParamStmt(convFunc string, param *model.Field) ast.Stmt {
	paramVar := utils.FirstToLower(param.Instance)
	return astDefineStmt(
		ast.NewIdent(paramVar),
		&ast.CallExpr{
			Fun: astSelectorExpr("utils", convFunc),
			Args: []ast.Expr{
				ast.NewIdent("webCtx"),
				astStringExpr(paramVar),
			},
		},
	)
}

func defineParamWithError(convFunc string, param *model.Field) []ast.Stmt {
	paramVar := utils.FirstToLower(param.Instance)
	return []ast.Stmt{
		astDefineStmtMany(
			[]ast.Expr{
				ast.NewIdent(paramVar),
				ast.NewIdent("err"),
			},
			&ast.CallExpr{
				Fun: astSelectorExpr("utils", convFunc),
				Args: []ast.Expr{
					ast.NewIdent("webCtx"),
					astStringExpr(paramVar),
				},
			},
		),
		errorReturnStmts(),
	}
}

func defineParamByParser(convFunc string, param *model.Field, packageName string) []ast.Stmt {
	paramVar := utils.FirstToLower(param.Instance)
	return []ast.Stmt{
		// [code] {{ParamInstance}} := &{{Package}}.{{ParamType}}{}
		astDefineStmt(
			ast.NewIdent(paramVar),
			astDeclareRef(
				param.Type,
				nil,
			),
		),
		// [code] err := utils.{{ParamSource}}Parser(webCtx, {{ParamInstance}})
		astAssignStmt(
			ast.NewIdent("err"),
			&ast.CallExpr{
				Fun: astSelectorExpr("utils", convFunc),
				Args: []ast.Expr{
					ast.NewIdent("webCtx"),
					ast.NewIdent(paramVar),
				},
			},
		),
		errorReturnStmts(),
	}
}
func genWebBodyParam(bodyParam *model.Field, packageName string, funcNode *model.Func) []ast.Stmt {
	if bodyParam != nil {
		bodyAstType := bodyParam.Type
		paramVar := utils.FirstToLower(bodyParam.Instance)
		switch bodyAstType.(type) {
		case *ast.ArrayType:
			byteArr := bodyAstType.(*ast.ArrayType)
			if byteArr.Elt.(*ast.Ident).String() == "byte" {
				return []ast.Stmt{
					astDefineStmt(
						ast.NewIdent(paramVar),
						&ast.CallExpr{
							Fun: astSelectorExpr("utils", "Body"),
							Args: []ast.Expr{
								ast.NewIdent("webCtx"),
							},
						},
					),
				}
			}
		case *ast.Ident:
			typeIdent := bodyAstType.(*ast.Ident).String()
			if typeIdent == "string" {
				return []ast.Stmt{
					astDefineStmt(
						ast.NewIdent(paramVar),
						&ast.CallExpr{
							Fun: astSelectorExpr("utils", "BodyString"),
							Args: []ast.Expr{
								ast.NewIdent("webCtx"),
							},
						},
					),
				}
			} else if utils.IsFirstLower(typeIdent) {
				utils.Failuref("%s %s, unsupport type %s", funcNode.Loc.String(), bodyParam.Comment, utils.TypeToString(bodyAstType))
			}
		}
		return defineParamByParser("BodyParser", bodyParam, packageName)
	}
	return nil
}

func genWebHeaderParams(headerParams []*model.Field, packageName string, funcNode *model.Func) []ast.Stmt {
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
						utils.Failuref("%s %s, unsupport type %s", funcNode.Loc.String(), param.Comment, utils.TypeToString(paramType))
					}
				}
			}
			stmts = append(stmts, defineParamByParser("HeaderParser", param, packageName)...)
		}
		return stmts
	}
	return nil
}
func genWebQueryParams(queryParams []*model.Field, packageName string, funcNode *model.Func) []ast.Stmt {
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
						utils.Failuref("%s %s, unsupport type %s", funcNode.Loc.String(), param.Comment, utils.TypeToString(paramType))
					}
				}
			}
			stmts = append(stmts, defineParamByParser("QueryParser", param, packageName)...)
		}
		return stmts
	}
	return nil
}

func genWebPathParams(pathParams []*model.Field, packageName string, funcNode *model.Func) []ast.Stmt {
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
						utils.Failuref("%s %s, unsupport type %s", funcNode.Loc.String(), param.Comment, utils.TypeToString(paramType))
					}
				}
			}
			stmts = append(stmts, defineParamByParser("ParamsParser", param, packageName)...)
		}
		return stmts
	}
	return nil
}
func genWebFormParams(formParams []*model.Field, packageName string, funcNode *model.Func) []ast.Stmt {
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
						utils.Failuref("%s %s, unsupport type %s", funcNode.Loc.String(), param.Comment, utils.TypeToString(paramType))
					}
				}
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
			}
			stmts = append(stmts, defineParamByParser("FormParser", param, packageName)...)
		}
		return stmts
	}
	return nil
}
