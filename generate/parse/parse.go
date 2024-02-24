package parse

import (
	. "github.com/ellisez/inject-golang/generate/global"
	"github.com/ellisez/inject-golang/generate/model"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

func annotateParse(text string) []string {
	prefixLen := len("// ")
	text = text[prefixLen:]
	text = strings.TrimSpace(text)
	return strings.Split(text, " ")
}

type Parser struct {
	Result *model.AnnotateInfo
}

func New() *Parser {
	result := &model.AnnotateInfo{}
	return &Parser{Result: result}
}

// DoParse
// 解析代码 -> ast
func (p *Parser) DoParse(filename string) error {
	ext := filepath.Ext(filename)
	if ext == ".go" {
		fileSet := token.NewFileSet()
		astFile, err := parser.ParseFile(fileSet, filename, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		// 解析根包名
		info := &model.PackageInfo{}
		info.Dirname = filepath.Dir(filename)
		info.Package = astFile.Name.String()

		if RootPackage == "" && RootDirectory == info.Dirname {
			RootPackage = info.Package
		}

		decls := astFile.Decls

		for _, decl := range decls {
			switch decl.(type) {
			case *ast.GenDecl:
				genDecl := decl.(*ast.GenDecl)
				switch genDecl.Tok {
				case token.VAR:
					break
				case token.TYPE:
					typeSpec := genDecl.Specs[0].(*ast.TypeSpec)
					switch typeSpec.Type.(type) {
					case *ast.InterfaceType:
						break
					case *ast.StructType:
						p.StructParse(genDecl, info)
						break
					}
					break
				default:
				}
				break
			case *ast.FuncDecl:
				funcDecl := decl.(*ast.FuncDecl)
				p.FuncParse(funcDecl, info)
				break
			}
		}

	}
	return nil
}
