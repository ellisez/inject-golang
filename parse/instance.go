package parse

import (
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
)

func (p *Parser) InstanceParse(funcDecl *ast.FuncDecl, commonFunc *model.CommonFunc, comments []*model.Comment) {
	instanceNode := model.NewProvide()
	instanceNode.Mode = "singleton"
	instanceNode.CommonFunc = commonFunc

	commonFunc.Loc = p.Ctx.FileSet.Position(funcDecl.Pos())
	for _, comment := range comments {
		args := comment.Args
		argsLen := len(args)
		if argsLen == 0 {
			continue
		}
		annotateName := args[0]
		switch annotateName {
		case "@provide":
			if argsLen < 2 {
				utils.Failuref(`%s %s, Instance must be specified`, commonFunc.Loc.String(), instanceNode.Comment)
			}
			instance := args[1]
			if utils.IsFirstLower(instance) {
				utils.Failuref(`%s %s, Instance "%s" must be capitalized with the first letter`, commonFunc.Loc.String(), instanceNode.Comment, instance)
			}
			instanceNode.Instance = instance

			if argsLen >= 3 {
				mode := args[2]
				if mode != "" && mode != "_" {
					instanceNode.Mode = mode
				}
			}
			instanceNode.Comment = comment.Comment
			break
		case "@order":
			if argsLen >= 2 {
				order := args[1]
				if order != "" && order != "_" {
					instanceNode.Order = order
				}
			}
			break
		case "@constructor":
			if argsLen < 2 {
				utils.Failuref(`%s %s, Constructor must be specified`, commonFunc.Loc.String(), instanceNode.Comment)
			}
			instanceNode.Constructor = args[1]
			break
		}
	}

	if funcDecl.Type.Results == nil {
		utils.Failuref("%s %s.%s() is not a valid constructor, missing return.", commonFunc.Loc.String(), commonFunc.Package, commonFunc.FuncName)
	}
	if len(funcDecl.Type.Results.List) > 1 {
		utils.Failuref("%s %s.%s() is not a valid constructor, It should only one return.", commonFunc.Loc.String(), commonFunc.Package, commonFunc.FuncName)
	}
	instanceNode.Type = funcDecl.Type.Results.List[0].Type

	switch instanceNode.Mode {
	case "singleton":
		p.Ctx.SingletonInstances = append(p.Ctx.SingletonInstances, instanceNode)
		break
	case "multiple":
		p.Ctx.MultipleInstances = append(p.Ctx.MultipleInstances, instanceNode)
		break
	default:
		utils.Failuref(`%s %s, Mode "%s" is invalid`, commonFunc.Loc.String(), instanceNode.Comment, instanceNode.Mode)
	}

}
