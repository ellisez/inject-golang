package utils

import (
	"go/ast"
	"regexp"
)

func TypeToString(astType ast.Expr) string {
	var str string
	switch astType.(type) {
	case *ast.Ident:
		str = astType.(*ast.Ident).String()
		break
	case *ast.StarExpr:
		starExpr := astType.(*ast.StarExpr)
		str = "*" + TypeToString(starExpr.X)
		break
	case *ast.SelectorExpr:
		selectorExpr := astType.(*ast.SelectorExpr)
		str = TypeToString(selectorExpr.X) + "." + selectorExpr.Sel.String()
		break
	case *ast.ChanType:
		chanType := astType.(*ast.ChanType)
		str = "chan " + TypeToString(chanType.Value)
		break
	}
	return str
}

func ShortType(typeString string) string {
	pattern := regexp.MustCompile(`(\*)?(?:(\w+)\.)?([A-Z]\w+)`)
	subMatch := pattern.FindStringSubmatch(typeString)
	return subMatch[2]
}

func TypeToAst(typeString string) (typeExpr ast.Expr) {
	pattern := regexp.MustCompile(`(\*)?(?:(\w+)\.)?([A-Z]\w+)`)
	subMatch := pattern.FindStringSubmatch(typeString)
	typeExpr = &ast.Ident{Name: subMatch[2]}
	if subMatch[1] != "" {
		typeExpr = &ast.SelectorExpr{
			X: typeExpr,
		}
	}
	if subMatch[0] != "" {
		typeExpr = &ast.StarExpr{
			X: typeExpr,
		}
	}
	return typeExpr
}
