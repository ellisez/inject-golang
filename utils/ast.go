package utils

import (
	"fmt"
	"github.com/ellisez/inject-golang/model"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
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
		break
	case *ast.StarExpr:
		starExpr := astType.(*ast.StarExpr)
		str = TypeShortName(starExpr.X)
		break
	case *ast.SelectorExpr:
		selectorExpr := astType.(*ast.SelectorExpr)
		str = selectorExpr.Sel.String()
		break
	case *ast.ChanType:
		chanType := astType.(*ast.ChanType)
		str = "chan" + FirstToUpper(TypeShortName(chanType.Value))
		break
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
		selectorExpr.X = selectorTypeWithNoPackage(subSelectorExpr, packageName)
		break
	case *ast.Ident:
		ident := selectorExpr.X.(*ast.Ident).String()
		if ident == packageName {
			return selectorExpr.Sel
		}
		break
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
		return astType
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
		starExpr.X = TypeWithPackage(starExpr.X, packageName)
		return starExpr
	case *ast.ChanType:
		chanType := astType.(*ast.ChanType)
		chanType.Value = TypeWithPackage(chanType.Value, packageName)
		return chanType
	case *ast.FuncType:
		return astType
	}
	return astType
}

func FieldName(field *model.Field) string {
	if field.Name != "" {
		return field.Name
	}
	return TypeShortName(field.Type)
}

func FindParam(funcInfo *model.Func, fieldName string) *model.Field {
	for _, field := range funcInfo.Params {
		if fieldName == field.Name {
			return field
		}
	}
	return nil
}
func ToFile(field *ast.Field) *model.Field {
	var fieldName string
	if field.Names != nil {
		fieldName = field.Names[0].String()
	}
	return &model.Field{
		Name: fieldName,
		Type: field.Type,
	}
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

func AddUniqueImport(imports []*ast.ImportSpec, importName string, importPath string) []*ast.ImportSpec {
	importPath = fmt.Sprintf(`"%s"`, importPath)
	var astImport *ast.ImportSpec
	for _, aImport := range imports {
		if importPath == aImport.Path.Value {
			astImport = aImport
			break
		}
	}
	if astImport == nil {
		astImport = &ast.ImportSpec{
			Name: &ast.Ident{Name: importName},
			Path: &ast.BasicLit{
				Value: importPath,
			},
		}

		return append(imports, astImport)
	}
	return imports
}

func UnusedImports(astFile *ast.File) {
	var imports []*ast.ImportSpec
	addImport := func(spec *ast.ImportSpec) {
		for _, importSpec := range imports {
			if importSpec.Path == spec.Path {
				return
			}
		}
		imports = append(imports, spec)
	}
	for _, spec := range astFile.Imports {
		for _, ident := range astFile.Unresolved {
			specName := spec.Name
			if specName != nil {
				if ident.String() == specName.String() {
					addImport(spec)
				}
			} else {
				specPath := spec.Path.Value
				importPath := specPath[1 : len(specPath)-1]
				if ok, _ := IsAllowedPackageName(importPath, ident.String()); ok {
					addImport(spec)
				}
			}
		}
	}
	astFile.Imports = imports

	// copy to genDecl
	specs := make([]ast.Spec, len(imports))
	for i, spec := range imports {
		specs[i] = spec
	}

	// find import genDecl
	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok == token.IMPORT {
			genDecl.Specs = specs
		}
		break
	}
}

func FieldSetter(fieldInfo *model.Field) string {
	fieldSetter := fieldInfo.Instance
	if fieldSetter == "" {
		fieldSetter = FieldName(fieldInfo)
	}
	return "Set" + fieldSetter
}

func FieldGetter(fieldInfo *model.Field) string {
	fieldGetter := fieldInfo.Instance
	if fieldGetter == "" {
		fieldGetter = FieldName(fieldInfo)
	}
	return fieldGetter
}
