package gen

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/token"
	"path/filepath"
)

func DoGen(ctx *model.Ctx) error {
	genDir := filepath.Join(Mod.Path, GenPackage)

	var err error
	if FlagAll || FlagSingleton {
		err = genSingletonFile(ctx, genDir)
		if err != nil {
			return err
		}
	}

	if FlagAll || FlagMultiple {
		err = genMultipleFile(ctx, genDir)
		if err != nil {
			return err
		}
	}

	if FlagAll || FlagFunc {
		err = genFuncFile(ctx, genDir)
		if err != nil {
			return err
		}

		err = genMethodFile(ctx, genDir)
		if err != nil {
			return err
		}
	}

	if FlagAll || FlagWeb {
		err = genWebUtilsFile(ctx, genDir)
		if err != nil {
			return err
		}

		err = genWebFile(ctx, genDir)
		if err != nil {
			return err
		}
	}

	err = genCtxFile(ctx, genDir)
	if err != nil {
		return err
	}

	err = genFactoryFile(ctx, genDir)
	if err != nil {
		return err
	}
	return nil
}

func astSelectorExpr(x string, sel string) *ast.SelectorExpr {
	selectorExpr := new(ast.SelectorExpr)
	selectorExpr.X = ast.NewIdent(x)
	selectorExpr.Sel = ast.NewIdent(sel)
	return selectorExpr
}

func astSelectorExprRecur(x ast.Expr, sel string) *ast.SelectorExpr {
	selectorExpr := new(ast.SelectorExpr)
	selectorExpr.X = x
	selectorExpr.Sel = ast.NewIdent(sel)
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
	var astName []*ast.Ident
	if name != "" {
		astName = []*ast.Ident{
			ast.NewIdent(name),
		}
	}
	return &ast.Field{
		Names: astName,
		Type:  typeExpr,
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
		Name: ast.NewIdent(name),
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
	var fieldList *ast.FieldList
	if fields != nil {
		fieldList = &ast.FieldList{
			List: fields,
		}
	}
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.StructType{
					Fields: fieldList,
				},
			},
		},
	}
}

func astCtxGetter(doc string, getter string, privateName string, fieldType ast.Expr) *ast.FuncDecl {
	astDoc := &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: doc,
		},
	}}

	getterDecl := astFuncDecl(
		[]*ast.Field{
			astField("ctx", astStarExpr(ast.NewIdent("Ctx"))),
		},
		getter,
		nil,
		[]*ast.Field{
			astField("", fieldType),
		},
		[]ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					astSelectorExpr("ctx", privateName),
				},
			},
		},
	)
	getterDecl.Doc = astDoc
	return getterDecl
}

func astCtxSetter(doc string, getter string, privateName string, fieldType ast.Expr) *ast.FuncDecl {
	astDoc := &ast.CommentGroup{List: []*ast.Comment{
		{
			Text: doc,
		},
	}}

	setterDecl := astFuncDecl(
		[]*ast.Field{
			astField("ctx", astStarExpr(ast.NewIdent("Ctx"))),
		},
		getter,
		[]*ast.Field{
			astField(privateName, fieldType),
		},
		nil,
		[]ast.Stmt{
			astAssignStmt(astSelectorExpr("ctx", privateName),
				ast.NewIdent(privateName),
			),
		},
	)
	setterDecl.Doc = astDoc
	return setterDecl
}

func astNewInstance(instance *model.Provide, ctxVar string) ast.Expr {
	instanceName := instance.Instance
	switch instance.Mode {
	case "singleton":
		instanceVar := utils.FirstToLower(instanceName)
		return astSelectorExpr(ctxVar, instanceVar)
	case "multiple":
		return &ast.CallExpr{
			Fun: astSelectorExpr(ctxVar, "New"+instanceName),
		}
	case "argument":
		return &ast.TypeAssertExpr{
			X: &ast.IndexExpr{
				X:     astSelectorExpr(ctxVar, ArgumentVar),
				Index: astStringExpr(instanceName),
			},
			Type: instance.Type,
		}

	}
	return nil
}

func astInstanceProxyFunc(instanceFunc *model.Func, instanceName string, params ...*ast.Field) *ast.FuncDecl {
	if instanceName == "" {
		instanceName = instanceFunc.FuncName
	}

	instanceVar := utils.FirstToLower(CtxType)

	params = append(params, astInstanceProxyParams(instanceFunc)...)
	var results []*ast.Field
	for _, result := range instanceFunc.Results {
		results = append(results, astField(result.Name, result.Type))
	}

	return astFuncDecl(
		[]*ast.Field{
			astField(instanceVar, astStarExpr(ast.NewIdent(CtxType))),
		},
		instanceName,
		params,
		results,
		nil,
	)
}
func astInstanceProxyParams(instanceFunc *model.Func) []*ast.Field {
	var params []*ast.Field
	for _, param := range instanceFunc.Params {
		if param.Source == "" {
			paramVar := utils.FieldVar(param)
			params = append(params, astField(paramVar, param.Type))
		}
	}
	return params
}

func astCallVarDefine(instanceFunc *model.Func, ctx *model.Ctx, ctxVar string) []ast.Stmt {
	var stmts []ast.Stmt
	for _, call := range instanceFunc.Calls {
		var lhs []ast.Expr

		fillLen := 0
		proxy := ctx.FuncOf(call.Instance)
		if proxy == nil {
			method := ctx.MethodOf(call.Instance)
			if method == nil {
				utils.Failuref(`%s %s, Instance "%s" is not found`, instanceFunc.Loc.String(), call.Comment, call.Instance)
			}
			if method.Type.Results == nil {
				utils.Failuref("%s %s, @injectCall must have at least one return value.", instanceFunc.Loc.String(), call.Comment)
			}

			fillLen = len(method.Type.Results.List) - len(call.Params)
			if fillLen < 0 {
				utils.Failuref(`%s %s, Params more than returns`, instanceFunc.Loc.String(), call.Comment)
			}
		} else {
			if proxy.Results == nil {
				utils.Failuref("%s %s, @injectCall must have at least one return value.", instanceFunc.Loc.String(), call.Comment)
			}

			fillLen = len(proxy.Results) - len(call.Params)
			if fillLen < 0 {
				utils.Failuref(`%s %s, Params more than returns`, instanceFunc.Loc.String(), call.Comment)
			}
		}

		for _, paramName := range call.Params {
			if paramName == "" || paramName == "_" {
				lhs = append(lhs, ast.NewIdent("_"))
				continue
			}
			varName := utils.NextVar(call.Instance + "." + paramName)
			lhs = append(lhs, ast.NewIdent(varName))
		}
		for i := 0; i < fillLen; i++ {
			lhs = append(lhs, ast.NewIdent("_"))
		}

		var argExpr ast.Expr
		argExpr = astSelectorExpr(ctxVar, call.Instance)
		argExpr = &ast.CallExpr{
			Fun: argExpr,
		}
		stmts = append(stmts, astDefineStmtMany(lhs, argExpr))
	}
	return stmts
}
func astInstanceCallExpr(handler ast.Expr, instanceFunc *model.Func, ctx *model.Ctx, ctxVar string) (*ast.CallExpr, []ast.Stmt) {
	stmts := astCallVarDefine(instanceFunc, ctx, ctxVar)
	var args []ast.Expr
	for _, param := range instanceFunc.Params {
		var argExpr ast.Expr
		switch param.Source {
		case "ctx":
			argExpr = ast.NewIdent(ctxVar)
		case "webCtx":
			argExpr = ast.NewIdent("webCtx")
		case "inject":
			paramInstance, _ := ctx.InstanceOf(param.Instance)
			if paramInstance == nil {
				utils.Failuref(`%s %s, Instance "%s" is not found`, param.Loc.String(), param.Comment, param.Instance)
			}
			argExpr = astNewInstance(paramInstance, ctxVar)
		case "func":
			proxy := ctx.FuncOf(param.Instance)
			if proxy == nil {
				method := ctx.MethodOf(param.Instance)
				if method == nil {
					utils.Failuref(`%s %s, Instance "%s" is not found`, param.Loc.String(), param.Comment, param.Instance)
				}
			}
			argExpr = astSelectorExpr(ctxVar, param.Instance)
		case "call":
			varName := utils.VarMap(param.Instance + "." + param.Name)
			if varName == "" {
				utils.Failuref(`%s %s, Instance "%s" is not found`, param.Loc.String(), param.Comment, param.Instance)
			}
			argExpr = ast.NewIdent(varName)
		case "":
			fallthrough
		default:
			paramVar := utils.FieldVar(param)
			argExpr = ast.NewIdent(paramVar)
		}

		switch param.Operator {
		case "&":
			argExpr = &ast.UnaryExpr{Op: token.AND, X: argExpr}
		case "*":
			argExpr = astStarExpr(argExpr)
		case "cast":
			argExpr = &ast.TypeAssertExpr{
				X:    argExpr,
				Type: param.Type,
			}
		}
		args = append(args, argExpr)
	}

	// [code] {{Package}}.{{Handler}}(...)
	return &ast.CallExpr{
		Fun:  handler,
		Args: args,
	}, stmts

}

func addImport(astFile *ast.File, ctx *model.Ctx, importName string, importPath string) error {
	var err error
	astFile.Imports, err = utils.AddUniqueImport(astFile.Imports, importName, importPath)
	if err != nil {
		return err
	}
	ctx.Imports, err = utils.AddUniqueImport(ctx.Imports, importName, importPath)
	if err != nil {
		return err
	}
	return nil
}

func astTypeToDeclare(typeExpr ast.Expr) ast.Expr {
	switch caseType := typeExpr.(type) {
	case *ast.StarExpr:
		return astDeclareRef(caseType.X, nil)
	case *ast.Ident:
		if utils.IsFirstLower(caseType.String()) {
			return nil
		}
	}
	return typeExpr
}

func addDecl(astFile *ast.File, genDecl ast.Decl) {
	if astFile.Decls == nil {
		astFile.Decls = make([]ast.Decl, 0)
	}
	astFile.Decls = append(astFile.Decls, genDecl)
}
