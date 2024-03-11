package parse

import (
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

func (p *Parser) CtxParse(funcDecl *ast.FuncDecl, funcInfo *model.FuncInfo) {
	if funcDecl.Doc != nil {
		ctxInfo := model.NewCtxInfoFromFunc(funcInfo)
		hasAnnotate := false
		for _, comment := range funcDecl.Doc.List {
			annotateArgs := annotateParse(comment.Text)
			argsLen := len(annotateArgs)
			if argsLen == 0 {
				continue
			}
			annotateName := annotateArgs[0]
			if annotateName == "@ctxProvide" {
				hasAnnotate = true
			} else if annotateName == "@proxy" {
				if argsLen >= 2 {
					proxy := annotateArgs[1]
					if proxy != "" && proxy != "_" {
						funcInfo.Proxy = proxy
					}
				} else {
					funcInfo.Proxy = funcInfo.FuncName
				}
				funcInfo.ProxyComment = comment.Text
			} else if annotateName == "@value" {
				if argsLen < 2 {
					utils.Failuref("%s, FieldName must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				fieldName := annotateArgs[1]
				if argsLen < 3 {
					utils.Failuref("%s, FieldType must be specified, at %s()", comment.Text, funcInfo.FuncName)
				}
				fieldType := annotateArgs[2]
				switch fieldType {
				case "string", "bool", "int", "float64", "uint":
				default:
					utils.Failuref("%s, FieldType unsupport \"%s\", at %s()", comment.Text, fieldType, funcInfo.FuncName)
				}

				value := &model.Value{Name: fieldName, Type: fieldType}
				ctxInfo.Values = append(ctxInfo.Values, value)
				if argsLen >= 4 {
					value.Default = annotateArgs[3]
				}
				value.Comment = comment.Text
			}
		}
		if hasAnnotate {
			p.Result.CtxInstances = append(p.Result.CtxInstances, ctxInfo)
		}

		if ctxInfo.Proxy != "" {
			addFuncOrMethodInstances(p.Result, funcInfo)
		}
	}
}
