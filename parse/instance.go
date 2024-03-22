package parse

import (
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"strings"
)

func (p *Parser) InstanceParse(funcDecl *ast.FuncDecl, commonFunc *model.CommonFunc, comments []*model.Comment) {
	instanceNode := model.NewProvide()
	instanceNode.CommonFunc = commonFunc

	commonFunc.Loc = p.Ctx.FileSet.Position(funcDecl.Pos())

	instanceValidate(instanceNode)
	for _, comment := range comments {
		args := comment.Args
		argsLen := len(args)
		if argsLen == 0 {
			continue
		}
		annotateName := args[0]
		switch annotateName {
		case "@provide":
			if argsLen >= 2 {
				instance := args[1]
				if instance != "" && instance != "_" {
					if utils.IsFirstLower(instance) {
						utils.Failuref(`%s %s, Instance "%s" must be capitalized with the first letter`, commonFunc.Loc.String(), instanceNode.Comment, instance)
					}
					instanceNode.Instance = instance
				}
			}

			if argsLen >= 3 {
				mode := args[2]
				if mode != "" && mode != "_" {
					switch mode {
					case "singleton", "multiple", "argument":
					default:
						utils.Failuref(`%s %s, Mode "%s" is invalid`, commonFunc.Loc.String(), instanceNode.Comment, instanceNode.Mode)
					}
					instanceNode.Mode = mode
				}
			}

			if argsLen >= 4 {
				typeStr := args[3]
				if typeStr != "" && typeStr != "_" {
					instanceNode.Type = utils.TypeToAst(typeStr)
				}
			}

			instanceNode.Comment = comment.Text
		case "@order":
			if argsLen >= 2 {
				order := args[1]
				if order != "" && order != "_" {
					instanceNode.Order = strings.TrimFunc(order, func(r rune) bool {
						return r == '"' || r == '`'
					})
				}
			}
		case "@handler":
			if argsLen < 2 {
				utils.Failuref(`%s %s, Handler must be specified`, commonFunc.Loc.String(), instanceNode.Comment)
			}
			instanceNode.Handler = args[1]
		}
	}

	addInstance(p.Ctx, instanceNode)
}

func instanceValidate(instance *model.Provide) {
	if len(instance.Results) > 1 {
		utils.Failuref("%s %s.%s() is not a valid constructor, It should only one return.", instance.Loc.String(), instance.Package, instance.FuncName)
	}
	instance.Type = instance.Results[0].Type
	instance.Instance = utils.TypeShortName(instance.Type)
}

func addInstance(ctx *model.Ctx, provide *model.Provide) {
	switch provide.Mode {
	case "singleton", "argument":
		instance, _ := ctx.SingletonOf(provide.Instance)
		if instance != nil {
			if !instance.Override {
				utils.Failuref(`%s %s, Instance "%s" Duplicate declaration`, provide.Loc.String(), provide.Comment, provide.Instance)
			}
			fmt.Printf(`Instance "%s" is Overrided by %s.%s`+"\n", provide.Instance, ImportAliasMap[provide.Package].Path, provide.FuncName)
			ctx.SingletonInstance.Replace(provide)
		} else {
			ctx.SingletonInstance.Add(provide)
		}
	case "multiple":
		instance := ctx.MultipleOf(provide.Instance)
		if instance != nil {
			if !instance.Override {
				utils.Failuref(`%s %s, Instance "%s" Duplicate declaration`, provide.Loc.String(), provide.Comment, provide.Instance)
			}
			fmt.Printf(`Instance "%s" is Overrided by %s.%s`+"\n", provide.Instance, ImportAliasMap[provide.Package].Path, provide.FuncName)
			ctx.MultipleInstance.Replace(provide)
			return
		} else {
			ctx.MultipleInstance.Add(provide)
		}
	}
}
