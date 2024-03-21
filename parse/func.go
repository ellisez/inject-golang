package parse

import (
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"strings"
)

func (p *Parser) FuncParse(funcDecl *ast.FuncDecl, currentImport *model.Import, currentImports []*ast.ImportSpec) {
	if funcDecl.Doc == nil {
		return
	}

	loc := p.Ctx.FileSet.Position(funcDecl.Pos())

	// main annotations
	var comments []*model.Comment
	var imports []*model.Import

	mode := "" // ""|@proxy|@provide|@webProvide|@middleware|@router
	for _, comment := range funcDecl.Doc.List {
		annotateArgs := annotateParse(comment.Text)
		argsLen := len(annotateArgs)
		if argsLen == 0 {
			continue
		}
		annotateName := annotateArgs[0]

		switch annotateName {
		case "@proxy", "@provide", "@webProvide", "@middleware", "@router":
			if mode != "" {
				utils.Failuref("%s %s, conflict with %s", loc.String(), comment.Text, mode)
			}
			mode = annotateName

			comments = append(comments, &model.Comment{
				Text: comment.Text,
				Args: annotateArgs,
			})
		case "@import":
			importInfo := &model.Import{}

			if argsLen < 2 {
				utils.Failuref("%s %s, Path must be specified", loc.String(), comment.Text)
			}
			importInfo.Path = annotateArgs[1]
			if argsLen >= 3 {
				importName := annotateArgs[2]
				if importName == "." {
					utils.Failuref("%s %s, Cannot support DotImport", loc.String(), comment.Text)
				}
				if importName == "" {
					importInfo.Package, _ = utils.GetPackageNameFromImport(importInfo.Path)
				} else {
					importInfo.Alias = importName
					importInfo.Package = importName
				}
			}

			importName, has := p.Ctx.ImportMap[importInfo.Path]
			if !has {
				p.Ctx.ImportMap[importInfo.Path] = importInfo.Alias
			} else if importName != importInfo.Alias {
				utils.Failuref(`%s %s, Alias conflicts with "%s"`, loc.String(), comment.Text, importName)
			}

			if importInfo.Path == currentImport.Path {
				currentImport.Alias = importInfo.Alias
			} else {
				imports = append(imports, importInfo)
			}
		default:
			comments = append(comments, &model.Comment{
				Text: comment.Text,
				Args: annotateArgs,
			})
		}
	}
	if mode == "" {
		return
	}

	targetImport := &model.Import{
		Package: GenPackage,
		Path:    Mod.Package + "/" + GenPackage,
	}

	/// parsing
	funcName := funcDecl.Name.String()

	funcNode := &model.Func{
		Package:  currentImport.Package,
		FuncName: funcName,
	}
	commonFunc := model.NewCommonFunc()
	commonFunc.Imports = append([]*model.Import{currentImport}, imports...)
	commonFunc.Func = funcNode
	astRecs := funcDecl.Recv

	if astRecs != nil {
		astRec := astRecs.List[0]
		rec, recImports := utils.ToFile(astRec, currentImport, currentImports, targetImport, imports)
		commonFunc.Imports = append(commonFunc.Imports, recImports...)
		rec.Loc = p.Ctx.FileSet.Position(astRecs.Pos())
		funcNode.Recv = rec
	}

	if funcDecl.Type.Params != nil {
		for _, astParam := range funcDecl.Type.Params.List {
			param, paramImports := utils.ToFile(astParam, currentImport, currentImports, targetImport, imports)
			commonFunc.Imports = append(commonFunc.Imports, paramImports...)
			param.Loc = p.Ctx.FileSet.Position(astParam.Pos())
			funcNode.Params = append(funcNode.Params, param)
		}
	}

	commonFunc.Loc = loc

	if funcDecl.Type.Results != nil {
		for _, asrResult := range funcDecl.Type.Results.List {
			result, resultImports := utils.ToFile(asrResult, currentImport, currentImports, targetImport, imports)
			commonFunc.Imports = append(commonFunc.Imports, resultImports...)
			result.Loc = p.Ctx.FileSet.Position(asrResult.Pos())
			funcNode.Results = append(funcNode.Results, result)
		}
	}

	var remainComments []*model.Comment
	for _, comment := range comments {
		annotateArgs := comment.Args
		argsLen := len(annotateArgs)
		if argsLen == 0 {
			continue
		}
		annotateName := annotateArgs[0]

		switch annotateName {
		case "@override":
			commonFunc.Override = true
		case "@injectParam":
			if argsLen < 2 {
				utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Text)
			}
			paramName := annotateArgs[1]
			param := utils.FindParam(funcNode, paramName)
			if param == nil {
				utils.Failuref("%s %s, ParamName not found", commonFunc.Loc.String(), comment.Text)
			}

			if argsLen >= 3 {
				paramInstance := annotateArgs[2]
				if paramInstance != "" && paramInstance != "_" {
					param.Instance = paramInstance
				}
			}

			if argsLen >= 4 {
				operator := annotateArgs[3]
				switch operator {
				case "", "&", "*", "cast":
					param.Operator = operator
				default:
					utils.Failuref(`%s %s, Operator "%s" not supported, only ["", "&", "*", "cast"] are allowed`, param.Loc.String(), comment.Text, operator)
				}
			}
			param.Comment = comment.Text
			param.Source = "inject"
		case "@injectFunc":
			if argsLen < 2 {
				utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Text)
			}
			paramName := annotateArgs[1]
			param := utils.FindParam(funcNode, paramName)
			if param == nil {
				utils.Failuref(`%s %s, ParamName "%s" not found`, commonFunc.Loc.String(), comment.Text, paramName)
			}

			if argsLen >= 3 {
				paramInstance := annotateArgs[2]
				if paramInstance != "" && paramInstance != "_" {
					param.Instance = paramInstance
				}
			}

			param.Comment = comment.Text
			param.Source = "func"
		case "@injectCall":
			if argsLen < 2 {
				utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Text)
			}
			paramNames := annotateArgs[1]
			if !strings.HasPrefix(paramNames, "[") || !strings.HasSuffix(paramNames, "]") {
				utils.Failuref("%s %s, ParamName must be wrapped in []", commonFunc.Loc.String(), comment.Text)
			}

			if argsLen < 3 {
				utils.Failuref("%s %s, Instance must be specified", commonFunc.Loc.String(), comment.Text)
			}
			instance := annotateArgs[2]
			call := &model.Call{
				Instance: instance,
			}

			paramNames = paramNames[1 : len(paramNames)-1]
			paramArr := strings.Split(paramNames, ",")
			for i, paramName := range paramArr {
				paramName = strings.TrimSpace(paramName)
				call.Params = append(call.Params, paramName)
				if paramName != "" && paramName != "_" {
					param := utils.FindParam(funcNode, paramName)
					if param == nil {
						utils.Failuref(`%s %s, ParamName "%s" not found`, commonFunc.Loc.String(), comment.Text, paramName)
					}

					param.Instance = instance
					param.Index = i

					param.Comment = comment.Text
					param.Source = "call"
				}
			}
			call.Comment = comment.Text
			funcNode.Calls = append(funcNode.Calls, call)
		case "@injectRecv":
			if argsLen < 2 {
				utils.Failuref("%s %s, RecvName must be specified", commonFunc.Loc.String(), comment.Text)
			}
			paramName := annotateArgs[1]
			if funcNode.Recv.Name != paramName {
				utils.Failuref("%s %s, RecvName not found", commonFunc.Loc.String(), comment.Text)
			}

			recv := funcNode.Recv
			if argsLen >= 3 {
				paramInstance := annotateArgs[2]
				if paramInstance != "" && paramInstance != "_" {
					recv.Instance = paramInstance
				}
			}
			recv.Comment = comment.Text
			recv.Source = "inject"
		case "@injectCtx":
			if argsLen < 2 {
				utils.Failuref("%s %s, ParamName must be specified", commonFunc.Loc.String(), comment.Text)
			}
			paramName := annotateArgs[1]
			param := utils.FindParam(funcNode, paramName)
			if param == nil {
				utils.Failuref("%s %s, ParamName not found", commonFunc.Loc.String(), comment.Text)
			}

			p.Ctx.InjectCtxMap[currentImport.Path] = append(p.Ctx.InjectCtxMap[currentImport.Path], param)
			param.Comment = comment.Text
			param.Source = "ctx"
		default:
			remainComments = append(remainComments, comment)
		}

	}

	switch mode {
	case "@proxy":
		p.ProxyParse(funcDecl, commonFunc, remainComments)
	case "@provide":
		p.InstanceParse(funcDecl, commonFunc, remainComments)
	case "@webProvide":
		p.WebParse(funcDecl, commonFunc, remainComments)
	case "@middleware":
		p.MiddlewareParse(funcDecl, commonFunc, remainComments)
	case "@router":
		p.RouterParse(funcDecl, commonFunc, remainComments)
	}
}
