package parse

import (
	. "github.com/ellisez/inject-golang/generate/global"
	"github.com/ellisez/inject-golang/generate/model"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"os"
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
	Result  *model.ModuleInfo
	FileSet *token.FileSet
}

func New() *Parser {
	return &Parser{
		Result:  &model.ModuleInfo{},
		FileSet: token.NewFileSet(),
	}
}

func (p *Parser) ModParse() error {
	filename := filepath.Join(p.Result.Dirname, "go.mod")
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	goDotMod, err := modfile.Parse("go.mod", bytes, nil)
	if err != nil {
		return err
	}
	p.Result.Mod = goDotMod.Module.Mod.Path
	return nil
}

// DoParse
// 解析代码 -> ast
func (p *Parser) DoParse(filename string) error {
	// exclude gen dir
	dirname := filepath.Dir(filename)
	if dirname == filepath.Join(p.Result.Dirname, GenPackage) {
		return nil
	}

	ext := filepath.Ext(filename)
	if ext == ".go" {
		astFile, err := parser.ParseFile(p.FileSet, filename, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		importPackage := p.Result.Mod
		if p.Result.Dirname != dirname {
			rel, err := filepath.Rel(p.Result.Dirname, dirname)
			if err != nil {
				return err
			}
			importPackage += "/" + filepath.ToSlash(rel)
		}

		// 解析包信息
		info := &model.PackageInfo{}
		info.Dirname = dirname
		info.Package = astFile.Name.String()
		info.Import = importPackage

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
