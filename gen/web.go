package gen

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func genWebFile(moduleInfo *model.ModuleInfo, dir string) error {
	fileDir := filepath.Join(dir, GenInternalPackage)
	filename := filepath.Join(fileDir, GenWebFilename)

	if len(moduleInfo.WebAppInstances) == 0 {
		err := os.Remove(filename)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		return nil
	}

	err := utils.CreateDirectoryIfNotExists(fileDir)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	if moduleInfo.WebAppInstances == nil {
		return nil
	}

	astFile := &ast.File{
		Name:  astIdent(GenPackage),
		Scope: ast.NewScope(nil),
	}

	genWebImportsAst(moduleInfo, astFile)

	genWebAppStartupAst(moduleInfo, astFile)

	genMiddlewareAst(moduleInfo, astFile)

	genRouterAst(moduleInfo, astFile)

	addFileDoc(astFile, "// Code generated by \"inject-golang -m web\"; DO NOT EDIT.")

	err = format.Node(file, token.NewFileSet(), astFile)
	if err != nil {
		return err
	}
	return nil
}

func genWebImportsAst(moduleInfo *model.ModuleInfo, astFile *ast.File) {
	uniqueImport(astFile, "", "fmt")
	uniqueImport(astFile, "", "github.com/gofiber/fiber/v2")
	uniqueImport(astFile, "", path.Join(Mod.Package, GenPackage, "utils"))

	for _, instance := range moduleInfo.WebAppInstances {
		uniqueImport(astFile, "", instance.Import)
		for _, importInfo := range instance.Imports {
			importName := importInfo.Name
			if importName == "_" {
				importName = ""
			}
			uniqueImport(astFile, importName, importInfo.Path)
			uniqueCtxImport(moduleInfo, importName, importInfo.Path)
		}

		for _, middleware := range instance.Middlewares {
			uniqueImport(astFile, "", middleware.Import)
			for _, importInfo := range middleware.Imports {
				importName := importInfo.Name
				if importName == "_" {
					importName = ""
				}
				uniqueImport(astFile, importName, importInfo.Path)
				uniqueCtxImport(moduleInfo, importName, importInfo.Path)
			}
		}

		for _, router := range instance.Routers {
			uniqueImport(astFile, "", router.Import)
			uniqueCtxImport(moduleInfo, "", router.Import)
			for _, importInfo := range router.Imports {
				importName := importInfo.Name
				if importName == "_" {
					importName = ""
				}
				uniqueImport(astFile, importName, importInfo.Path)
				uniqueCtxImport(moduleInfo, importName, importInfo.Path)
			}
		}
	}
	addImportDecl(astFile)
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
	recvVar := utils.FirstToLower(StructName)

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
								GenPackage,
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
			// [code] ctx.{{WebApp}}().Static({{Path}}, {{Path}}, {{...}})
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
						&ast.CallExpr{
							Fun: astSelectorExpr(recvVar, instance.WebApp),
						},
						"Static"),
					Args: args,
				},
			})
		}

		for _, middleware := range instance.Middlewares {
			// [code] ctx.{{WebApp}}().Group({{Path}}, {{Proxy}})
			stmts = append(stmts, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: astSelectorExprRecur(
						&ast.CallExpr{
							Fun: astSelectorExpr(recvVar, instance.WebApp),
						},
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
				// [code] ctx.{{WebApp}}().{{Method}}({{Path}}, {{Proxy}})
				stmts = append(stmts, &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: astSelectorExprRecur(
							&ast.CallExpr{
								Fun: astSelectorExpr(recvVar, instance.WebApp),
							},
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

				switch paramInfo.Source {
				case "ctx":
					// [code] ctx,
					args = append(args, astIdent(recvVar))
					break
				case "inject":
					injectMode := ""
					if moduleInfo.HasSingleton(paramInstance) {
						// [code] ctx.{{ParamInstance}}(),
						args = append(args,
							&ast.CallExpr{
								Fun: astSelectorExpr(recvVar, paramInstance),
							})
						injectMode = "singleton"
					}

					if injectMode == "" {
						if moduleInfo.HasWebApp(paramInstance) {
							// [code] ctx.{{ParamInstance}}(),
							args = append(args,
								&ast.CallExpr{
									Fun: astSelectorExpr(recvVar, paramInstance),
								})
							injectMode = "webApp"
						}
					}

					if injectMode == "" {
						if moduleInfo.HasMultiple(paramInstance) {
							// [code] ctx.New{{ParamInstance}}(),
							args = append(args, &ast.CallExpr{
								Fun: astSelectorExpr(recvVar, "New"+paramInstance),
							})
							injectMode = "multiple"
						}
					}

					if injectMode == "" {
						utils.Failuref("%s, \"%s\" No matching Instance, at %s()", paramInfo.Comment, paramInstance, instance.FuncName)
					}
					break
				default:
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
						&ast.CallExpr{
							Fun: astSelectorExpr(recvVar, instance.WebApp),
						},
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

		genDoc := &ast.Comment{
			Text: fmt.Sprintf("// Generate by annotations from %s.%s", instance.Package, instance.FuncName),
		}
		results := []*ast.Field{
			astField("", astIdent("error")),
		}
		funcDecl := astFuncDecl(
			[]*ast.Field{
				astField(recvVar, astStarExpr(astIdent(StructName))),
			},
			instance.Proxy,
			params,
			results,
			stmts,
		)
		funcDecl.Doc = &ast.CommentGroup{List: []*ast.Comment{
			{
				Text: fmt.Sprintf("// %s", instance.Proxy),
			},
			genDoc,
		}}
		addDecl(astFile, funcDecl)
		/// interface method field
		methodField := astField(
			instance.Proxy,
			&ast.FuncType{
				Params:  &ast.FieldList{List: params},
				Results: &ast.FieldList{List: results},
			},
		)
		methodField.Comment = &ast.CommentGroup{List: []*ast.Comment{
			genDoc,
		}}
		moduleInfo.CtxMethodFields = append(moduleInfo.CtxMethodFields, methodField)
	}

}

func defineParamStmt(convFunc string, param *model.FieldInfo) ast.Stmt {
	return astDefineStmt(
		astIdent(param.Instance),
		&ast.CallExpr{
			Fun: astSelectorExpr("utils", convFunc),
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
				Fun: astSelectorExpr("utils", convFunc),
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
					GenPackage,
				),
				nil,
			),
		),
		// [code] err := utils.{{ParamSource}}Parser(webCtx, {{ParamInstance}})
		astAssignStmt(
			astIdent("err"),
			&ast.CallExpr{
				Fun: astSelectorExpr("utils", convFunc),
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
							Fun: astSelectorExpr("utils", "Body"),
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
							Fun: astSelectorExpr("utils", "BodyString"),
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
