package parse

import (
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

func (p *Parser) FuncParse(funcDecl *ast.FuncDecl, packageName string, importPath string) {
	funcName := funcDecl.Name.String()

	funcNode := &model.Func{
		Package:  packageName,
		FuncName: funcName,
	}
	commonFunc := model.NewCommonFunc()
	commonFunc.Func = funcNode
	commonFunc.Imports = append(commonFunc.Imports, &model.Import{Path: importPath})

	astRec := funcDecl.Recv

	if astRec != nil {
		fieldRec := astRec.List[0]
		param := utils.ToFile(fieldRec, packageName, GenPackage)
		funcNode.Recv = param
		//funcNode.Params = append(funcNode.Params, param)
	}

	fillEmptyParam(funcDecl.Type, funcNode)

	commonFunc.Loc = p.Ctx.FileSet.Position(funcDecl.Pos())

	if funcDecl.Type.Results != nil {
		for _, result := range funcDecl.Type.Results.List {
			funcNode.Results = append(funcNode.Results, utils.ToFile(result, packageName, GenPackage))
		}
	}

	var comments []*model.Comment
	mode := "" // ""|@proxy|@provide|@webProvide|@middleware|@router
	if funcDecl.Doc != nil {
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
					utils.Failuref("%s %s, conflict with %s", commonFunc.Loc.String(), comment.Text, mode)
				}
				mode = annotateName

				comments = append(comments, &model.Comment{
					Comment: comment.Text,
					Args:    annotateArgs,
				})
			case "@import":
				importInfo := &model.Import{}
				commonFunc.Imports = append(commonFunc.Imports, importInfo)

				if argsLen < 2 {
					utils.Failuref("%s %s, Path must be specified", commonFunc.Loc.String(), comment.Text)
				}
				importInfo.Path = annotateArgs[1]
				if argsLen >= 3 {
					importName := annotateArgs[2]
					if importName == "." {
						utils.Failuref("%s %s, Cannot support DotImport", commonFunc.Loc.String(), comment.Text)
					}
					if importName != "" {
						importInfo.Name = importName
					}
				}
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
					pointer := annotateArgs[3]
					switch pointer {
					case "", "&", "*":
						param.Pointer = pointer
					default:
						utils.Failuref(`%s %s, Pointer "%s" not supported, only ["", "&", "*"] are allowed`, commonFunc.Loc.String(), comment.Text, pointer)
					}
				}
				param.Comment = comment.Text
				param.Source = "inject"
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
				paramInfo := utils.FindParam(funcNode, paramName)
				if paramInfo == nil {
					utils.Failuref("%s %s, ParamName not found", commonFunc.Loc.String(), comment.Text)
				}

				paramInfo.Comment = comment.Text
				paramInfo.Source = "ctx"
			default:
				comments = append(comments, &model.Comment{
					Comment: comment.Text,
					Args:    annotateArgs,
				})
			}

		}
	}

	switch mode {
	case "@proxy":
		p.ProxyParse(funcDecl, commonFunc, comments)
	case "@provide":
		p.InstanceParse(funcDecl, commonFunc, comments)
	case "@webProvide":
		p.WebParse(funcDecl, commonFunc, comments)
	case "@middleware":
		p.MiddlewareParse(funcDecl, commonFunc, comments)
	case "@router":
		p.RouterParse(funcDecl, commonFunc, comments)
	}
}
func fillEmptyParam(funcType *ast.FuncType, funcNode *model.Func) {
	for _, field := range funcType.Params.List {
		funcNode.Params = append(funcNode.Params, utils.ToFile(field, funcNode.Package, GenPackage))
	}
}
