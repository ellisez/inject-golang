package utils

import (
	"fmt"
	"github.com/ellisez/inject-golang/model"
	"go/ast"
	"regexp"
	"strings"
)

func TypeToString(astType ast.Expr) string {
	var str string
	switch astType.(type) {
	case *ast.Ident:
		str = astType.(*ast.Ident).String()
	case *ast.StarExpr:
		starExpr := astType.(*ast.StarExpr)
		str = "*" + TypeToString(starExpr.X)
	case *ast.SelectorExpr:
		selectorExpr := astType.(*ast.SelectorExpr)
		str = TypeToString(selectorExpr.X) + "." + selectorExpr.Sel.String()
	case *ast.ChanType:
		chanType := astType.(*ast.ChanType)
		str = "chan " + TypeToString(chanType.Value)
	case *ast.FuncType:
		funcType := astType.(*ast.FuncType)
		params := ""
		if funcType.Params != nil {
			for i, field := range funcType.Params.List {
				if i > 0 {
					params += ","
				}
				paramType := TypeToString(field.Type)
				if field.Names != nil {
					params += field.Names[0].String() + " " + paramType
				} else {
					params += paramType
				}

			}
		}
		results := ""
		if funcType.Results != nil {
			rl := len(funcType.Results.List)
			for i, field := range funcType.Results.List {
				if i > 0 {
					results += ","
				}
				resultType := TypeToString(field.Type)
				if field.Names != nil {
					results += field.Names[0].String() + " " + resultType
				} else {
					results += resultType
				}
			}
			if rl > 1 {
				results = "(" + results + ")"
			}
		}
		str = fmt.Sprintf("func(%s) %s", params, results)
	}
	return str
}

func TypeShortName(astType ast.Expr) string {
	var str string
	switch astType.(type) {
	case *ast.Ident:
		str = astType.(*ast.Ident).String()
	case *ast.StarExpr:
		starExpr := astType.(*ast.StarExpr)
		str = TypeShortName(starExpr.X)
	case *ast.SelectorExpr:
		selectorExpr := astType.(*ast.SelectorExpr)
		str = selectorExpr.Sel.String()
	case *ast.ChanType:
		chanType := astType.(*ast.ChanType)
		str = "chan" + FirstToUpper(TypeShortName(chanType.Value))
	case *ast.FuncType:
		str = ""
	}
	return str
}

func TypeToAst(typeString string) (typeExpr ast.Expr) {
	pattern := regexp.MustCompile(`(\*)?(?:(\w+)\.)?(\w+)`)
	subMatch := pattern.FindStringSubmatch(typeString)

	sel := &ast.Ident{Name: subMatch[3]}
	typeExpr = sel
	if subMatch[2] != "" {
		typeExpr = &ast.SelectorExpr{
			X:   &ast.Ident{Name: subMatch[2]},
			Sel: sel,
		}
	}
	if subMatch[1] != "" {
		typeExpr = &ast.StarExpr{
			X: typeExpr,
		}
	}
	return typeExpr
}

func AccessType(astType ast.Expr, definePackage string, accessPackage string) ast.Expr {
	if definePackage == accessPackage {
		return TypeWithNoPackage(astType, definePackage)
	} else {
		return TypeWithPackage(astType, definePackage)
	}

}
func selectorTypeWithNoPackage(selectorExpr *ast.SelectorExpr, packageName string) ast.Expr {
	switch selectorExpr.X.(type) {
	case *ast.SelectorExpr:
		subSelectorExpr := selectorExpr.X.(*ast.SelectorExpr)
		newSelectorExpr := *subSelectorExpr
		selectorExpr.X = selectorTypeWithNoPackage(&newSelectorExpr, packageName)
		return &newSelectorExpr
	case *ast.Ident:
		ident := selectorExpr.X.(*ast.Ident).String()
		if ident == packageName {
			return selectorExpr.Sel
		}
	}
	return selectorExpr
}
func TypeWithNoPackage(astType ast.Expr, packageName string) ast.Expr {
	switch astType.(type) {
	case *ast.SelectorExpr:
		selectorExpr := astType.(*ast.SelectorExpr)
		return selectorTypeWithNoPackage(selectorExpr, packageName)
	case *ast.Ident:
		return astType
	case *ast.StarExpr:
		starExpr := astType.(*ast.StarExpr)
		starExpr.X = TypeWithNoPackage(starExpr.X, packageName)
		return starExpr
	case *ast.ChanType:
		chanType := astType.(*ast.ChanType)
		chanType.Value = TypeWithNoPackage(chanType.Value, packageName)
		return chanType
	case *ast.FuncType:
		funcType := astType.(*ast.FuncType)
		newFuncType := *funcType
		if newFuncType.Params != nil {
			var params []*ast.Field
			for _, param := range newFuncType.Params.List {
				newParam := *param
				newParam.Type = TypeWithNoPackage(param.Type, packageName)
				params = append(params, &newParam)
			}
			newFuncType.Params = &ast.FieldList{List: params}
		}
		if newFuncType.Results != nil {
			var results []*ast.Field
			for _, result := range newFuncType.Results.List {
				newParam := *result
				newParam.Type = TypeWithNoPackage(result.Type, packageName)
				results = append(results, &newParam)
			}
			newFuncType.Results = &ast.FieldList{List: results}
		}
		return &newFuncType
	}
	return astType
}

func TypeWithPackage(astType ast.Expr, packageName string) ast.Expr {
	switch astType.(type) {
	case *ast.SelectorExpr:
		return astType
	case *ast.Ident:
		ident := astType.(*ast.Ident)
		if IsFirstUpper(ident.String()) {
			return &ast.SelectorExpr{
				X:   &ast.Ident{Name: packageName},
				Sel: ident,
			}
		} else {
			return ident
		}
	case *ast.StarExpr:
		starExpr := astType.(*ast.StarExpr)
		newStarExpr := *starExpr
		newStarExpr.X = TypeWithPackage(newStarExpr.X, packageName)
		return &newStarExpr
	case *ast.ChanType:
		chanType := astType.(*ast.ChanType)
		newChanType := *chanType
		newChanType.Value = TypeWithPackage(newChanType.Value, packageName)
		return &newChanType
	case *ast.FuncType:
		funcType := astType.(*ast.FuncType)
		newFuncType := *funcType
		if newFuncType.Params != nil {
			var params []*ast.Field
			for _, param := range newFuncType.Params.List {
				newParam := *param
				newParam.Type = TypeWithPackage(param.Type, packageName)
				params = append(params, &newParam)
			}
			newFuncType.Params = &ast.FieldList{List: params}
		}
		if newFuncType.Results != nil {
			var results []*ast.Field
			for _, result := range newFuncType.Results.List {
				newParam := *result
				newParam.Type = TypeWithPackage(result.Type, packageName)
				results = append(results, &newParam)
			}
			newFuncType.Results = &ast.FieldList{List: results}
		}
		return &newFuncType
	}
	return astType
}

func FieldName(field *model.Field) string {
	if field.Name != "" {
		return field.Name
	}
	return TypeShortName(field.Type)
}

func FieldVar(field *model.Field) string {
	fieldVar := field.Name
	if fieldVar == "" {
		fieldVar = field.Instance
	}
	return FirstToLower(fieldVar)
}

func FindParam(funcInfo *model.Func, fieldName string) *model.Field {
	for _, field := range funcInfo.Params {
		if fieldName == field.Name {
			return field
		}
	}
	return nil
}
func ToFile(field *ast.Field, definePackage string, accessPackage string) *model.Field {
	var fieldName string
	if field.Names != nil {
		fieldName = field.Names[0].String()
	}
	t := AccessType(field.Type, definePackage, accessPackage)
	f := &model.Field{
		Package: definePackage,
		Name:    fieldName,
		Type:    t,
	}
	f.Instance = FirstToUpper(FieldName(f))
	return f
}

func IsBasicType(typeStr string) bool {
	if strings.HasPrefix(typeStr, "*") {
		return false
	}
	if strings.Contains(typeStr, ".") {
		return false
	}
	if IsFirstUpper(typeStr) {
		return false
	}
	return true
}

func IsBasicAstType(typeExpr ast.Expr) bool {
	ident, ok := typeExpr.(*ast.Ident)
	if !ok {
		return false
	}
	if IsFirstUpper(ident.String()) {
		return false
	}
	return true
}

func AddUniqueImport(imports []*ast.ImportSpec, importName string, importPath string) ([]*ast.ImportSpec, error) {
	relImport := importName
	if relImport == "" || relImport == "_" {
		relImport, _ = GetPackageNameFromImport(importPath)
	}
	importPathValue := fmt.Sprintf(`"%s"`, importPath)
	var astImport *ast.ImportSpec
	for _, aImport := range imports {
		aImportPath := aImport.Path.Value
		aImportPath = aImportPath[1 : len(aImportPath)-1]
		aImportName := ""
		if aImport.Name != nil {
			aImportName = aImport.Name.String()
		}
		relAImport := aImportName
		if relAImport == "" {
			relAImport, _ = GetPackageNameFromImport(aImportPath)
		}

		if importPath == aImportPath {
			if relImport != relAImport {
				return nil, fmt.Errorf(`@import "%s" aliases "%s" and "%s" conflict`, importPath, aImportName, importName)
			}
			astImport = aImport
			break
		} else {
			if relImport == relAImport {
				return nil, fmt.Errorf(`@import %s "%s" conflicts with @import %s "%s", try to change alias`, importPath, importName, aImportPath, importName)
			}
		}
	}
	if astImport == nil {
		astImport = &ast.ImportSpec{
			Name: &ast.Ident{Name: importName},
			Path: &ast.BasicLit{
				Value: importPathValue,
			},
		}

		return append(imports, astImport), nil
	}
	return imports, nil
}

func StringLit(lit *ast.BasicLit) string {
	return strings.TrimFunc(lit.Value, func(r rune) bool {
		return r == '"' || r == '`'
	})
}
