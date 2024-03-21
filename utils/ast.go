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

func TypeForChangePackage(astType ast.Expr, currentImport *model.Import, currentImports []*ast.ImportSpec, changeImport *model.Import, changeImports []*model.Import) (ast.Expr, []*model.Import) {
	var imports []*model.Import

	if currentImport.Alias == changeImport.Alias && currentImport.Path == changeImport.Path {
		return astType, nil
	}
	switch astType.(type) {
	case *ast.SelectorExpr: // ?.yyy
		selectorExpr := astType.(*ast.SelectorExpr)
		ident, ok := selectorExpr.X.(*ast.Ident)
		if ok {
			var matchImport *model.Import
			var matchRelName string
			for _, astImport := range currentImports {
				name := ""
				if astImport.Name != nil {
					name = astImport.Name.String()
				}
				matchRelName = RelPackageNameOfAst(astImport)
				if matchRelName == ident.String() {
					matchImport = &model.Import{
						Alias:   name,
						Package: matchRelName,
						Path:    ImportPathOf(astImport),
					}
					break
				}
			}
			if matchImport == nil {
				Failuref(`%s, Missing package "%s"`, currentImport.Path, ident.String())
			}
			for _, importNode := range changeImports {
				if importNode.Path == matchImport.Path {
					if importNode.Alias != matchImport.Alias {
						return &ast.SelectorExpr{
							X:   ast.NewIdent(importNode.Alias),
							Sel: selectorExpr.Sel,
						}, imports
					} else {
						return selectorExpr, imports
					}
				}
			}
			imports = append(imports, matchImport)
			return selectorExpr, imports
		}
		newX, newImports := TypeForChangePackage(selectorExpr.X, currentImport, currentImports, changeImport, changeImports)
		if newX != selectorExpr.X {
			selectorExpr = &ast.SelectorExpr{
				X:   newX,
				Sel: selectorExpr.Sel,
			}
		}
		imports = append(imports, newImports...)
		return selectorExpr, imports
	case *ast.Ident: // xxx
		ident := astType.(*ast.Ident)
		if IsBasicType(ident.String()) {
			return ident, imports
		}
		return &ast.SelectorExpr{
			X:   ast.NewIdent(currentImport.Package),
			Sel: ident,
		}, imports
	case *ast.StarExpr: // *?
		starExpr := astType.(*ast.StarExpr)
		newX, newImports := TypeForChangePackage(starExpr.X, currentImport, currentImports, changeImport, changeImports)
		if newX != starExpr.X {
			starExpr = &ast.StarExpr{
				X: newX,
			}
		}
		imports = append(imports, newImports...)
		return starExpr, imports
	case *ast.ArrayType: // []?
		arrayType := astType.(*ast.ArrayType)
		newElt, newImports := TypeForChangePackage(arrayType.Elt, currentImport, currentImports, changeImport, changeImports)
		if newElt != arrayType.Elt {
			arrayType = &ast.ArrayType{
				Elt: newElt,
			}
		}
		imports = append(imports, newImports...)
		return arrayType, imports
	case *ast.MapType: // map[?]?
		mapType := astType.(*ast.MapType)

		hasChange := false
		newMapType := &ast.MapType{}
		newKey, newKeyImports := TypeForChangePackage(mapType.Key, currentImport, currentImports, changeImport, changeImports)
		if newKey != mapType.Key {
			hasChange = true
		}
		imports = append(imports, newKeyImports...)
		newMapType.Key = newKey

		newValue, newValueImports := TypeForChangePackage(mapType.Value, currentImport, currentImports, changeImport, changeImports)
		if newValue != mapType.Value {
			hasChange = true
		}
		imports = append(imports, newValueImports...)
		newMapType.Value = newValue

		if hasChange {
			return newMapType, imports
		}
		return mapType, imports
	case *ast.FuncType: // func(p ?) (r ?)
		funcType := astType.(*ast.FuncType)

		hasChange := false
		newFuncType := &ast.FuncType{
			TypeParams: funcType.TypeParams,
		}
		if funcType.Params != nil {
			newFuncType.Params = &ast.FieldList{}
			for _, field := range funcType.Params.List {
				newFieldType, newFieldImports := TypeForChangePackage(field.Type, currentImport, currentImports, changeImport, changeImports)
				if newFieldType != field.Type {
					hasChange = true
				}
				newFuncType.Params.List = append(newFuncType.Params.List, &ast.Field{
					Names: field.Names,
					Type:  newFieldType,
				})
				imports = append(imports, newFieldImports...)
			}

		}
		if funcType.Results != nil {
			newFuncType.Results = &ast.FieldList{}
			for _, field := range funcType.Results.List {
				newFieldType, newFieldImports := TypeForChangePackage(field.Type, currentImport, currentImports, changeImport, changeImports)
				if newFieldType != field.Type {
					hasChange = true
				}
				newFuncType.Results.List = append(newFuncType.Results.List, &ast.Field{
					Names: field.Names,
					Type:  newFieldType,
				})
				imports = append(imports, newFieldImports...)
			}

		}

		if hasChange {
			return newFuncType, imports
		}
		return funcType, imports
	}
	return nil, nil
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
func ToFile(field *ast.Field, currentImport *model.Import, currentImports []*ast.ImportSpec, changeImport *model.Import, changeImports []*model.Import) (*model.Field, []*model.Import) {
	var fieldName string
	if field.Names != nil {
		fieldName = field.Names[0].String()
	}
	//t := AccessType(field.Type, currentPackage, accessPackage, astImports)
	t, imports := TypeForChangePackage(field.Type, currentImport, currentImports, changeImport, changeImports)
	f := &model.Field{
		Package: currentImport.Package,
		Name:    fieldName,
		Type:    t,
	}
	f.Instance = FirstToUpper(FieldName(f))
	return f, imports
}

func IsBasicType(typeStr string) bool {
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
				return nil, fmt.Errorf(`@import "%s" aliases "%s" conflicts with "%s"`, importPath, aImportName, importName)
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
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: importPathValue,
			},
		}
		if importName != "" {
			astImport.Name = &ast.Ident{Name: importName}
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
