package parse

import (
	"fmt"
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
					case "singleton", "multiple":
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

			instanceNode.Comment = comment.Comment
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

	if instanceOverride(p.Ctx, instanceNode) {
		return
	}

	switch instanceNode.Mode {
	case "singleton":
		p.Ctx.SingletonInstances = append(p.Ctx.SingletonInstances, instanceNode)
	case "multiple":
		p.Ctx.MultipleInstances = append(p.Ctx.MultipleInstances, instanceNode)
	}
}

func instanceValidate(instance *model.Provide) {
	if len(instance.Results) > 1 {
		utils.Failuref("%s %s.%s() is not a valid constructor, It should only one return.", instance.Loc.String(), instance.Package, instance.FuncName)
	}
	instance.Type = instance.Results[0].Type
	instance.Instance = utils.TypeShortName(instance.Type)
}

func instanceOverride(ctx *model.Ctx, provide *model.Provide) bool {
	for i, instance := range ctx.SingletonInstances {
		if instance.GetInstance() == provide.Instance {
			if instance.GetOverride() {
				if provide.Mode == "singleton" {
					ctx.SingletonInstances[i] = provide
				} else {
					ctx.MultipleInstances = append(ctx.MultipleInstances, provide)
				}
				fmt.Printf(`Instance "%s" is Overrided by %s.%s`+"\n", provide.Instance, provide.Package, provide.FuncName)
				return true
			} else {
				utils.Failuref(`%s %s, Instance "%s" Duplicate declaration`, provide.Loc.String(), provide.Comment, provide.Instance)
			}
		}
	}
	for i, instance := range ctx.MultipleInstances {
		if instance.GetInstance() == provide.Instance {
			if instance.GetOverride() {
				if provide.Mode == "singleton" {
					ctx.SingletonInstances = append(ctx.SingletonInstances, provide)
				} else {
					ctx.MultipleInstances[i] = provide
				}
				fmt.Printf(`Instance "%s" is Overrided by %s.%s`+"\n", provide.Instance, provide.Package, provide.FuncName)
				return true
			} else {
				utils.Failuref(`%s %s, Instance "%s" Duplicate declaration`, provide.Loc.String(), provide.Comment, provide.Instance)
			}
		}
	}
	return false
}
