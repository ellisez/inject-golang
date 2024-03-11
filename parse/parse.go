package parse

import (
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
)

func annotateParse(text string) []string {
	prefixLen := len("// ")
	if len(text) < prefixLen {
		return nil
	}
	text = text[prefixLen:]

	var strs []string
	mode := "seek" // seek | endOfSpace | endOfDoubleQuote | endOfBackQuote
	lastIndex := 0

	add := func(index int) {
		if mode != "seek" {
			strs = append(strs, text[lastIndex:index])
			mode = "seek"
		}
		lastIndex = index
	}

	for index, char := range text {
		switch mode {
		case "seek":
			switch char {
			case ' ', '\t':
				continue
			case '"':
				mode = "endOfDoubleQuote"
				lastIndex = index
				break
			case '`':
				mode = "endOfBackQuote"
				lastIndex = index
				break
			default:
				mode = "endOfSpace"
				lastIndex = index
			}
		case "endOfSpace":
			switch char {
			case ' ', '\t':
				add(index)
				break
			default:

			}
			break
		case "endOfDoubleQuote":
			switch char {
			case '"':
				add(index + 1)
				break
			default:

			}
			break
		case "endOfBackQuote":
			switch char {
			case '`':
				add(index + 1)
				break
			default:

			}
			break
		}
	}
	if mode != "seek" {
		strs = append(strs, text[lastIndex:])
	}
	return strs
}

type Parser struct {
	*model.Mod
	Result *model.ModuleInfo
}

func ModParse(dirname string) (*model.Mod, error) {
	dirname = utils.JoinPath(dirname)
	mod := CacheModMap[dirname]
	if mod != nil {
		return mod, nil
	}

	mod = &model.Mod{}

	/// go.mod
	filename := filepath.Join(dirname, "go.mod")
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	goDotMod, err := modfile.Parse("go.mod", bytes, nil)
	if err != nil {
		return nil, err
	}
	mod.Path = dirname
	mod.Package = goDotMod.Module.Mod.Path
	mod.Version = goDotMod.Module.Mod.Version

	mod.Require = make(map[string]string)
	for _, r := range goDotMod.Require {
		mod.Require[r.Mod.Path] = r.Mod.Version
	}

	if Mod == nil {
		Mod = mod
		/// go.work
		filename = filepath.Join(dirname, "go.work")
		existsWork, err := utils.ExistsFile(filename)
		if err != nil {
			return nil, err
		}

		if existsWork {
			bytes, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			goDotWork, err := modfile.ParseWork("go.work", bytes, nil)
			if err != nil {
				return nil, err
			}
			mod.Work = make(map[string]string)
			for _, use := range goDotWork.Use {
				if use.Path == "." {
					mod.Work[mod.Package] = mod.Path
					continue
				}
				localMod, err := ModParse(use.Path)
				if err != nil {
					return nil, err
				}
				mod.Work[localMod.Package] = utils.JoinPath(use.Path)
			}
		}
	}
	CacheModMap[mod.Path] = mod

	return mod, nil
}

// DoParse
// 解析代码 -> ast
func (p *Parser) DoParse(filename string) error {
	// exclude gen dir
	dirname := filepath.Dir(filename)
	if dirname == filepath.Join(p.Path, GenPackage) {
		return nil
	}

	ext := filepath.Ext(filename)
	if ext == ".go" {
		astFile, err := parser.ParseFile(p.Result.FileSet, filename, nil, parser.ParseComments)
		if err != nil {
			utils.Failure(err.Error())
		}

		importPackage := p.Package
		if p.Path != dirname {
			rel, err := filepath.Rel(p.Path, dirname)
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

		isAllow, packageType := utils.IsAllowedPackageName(info.Import, info.Package)
		if !isAllow {
			utils.Failuref("Detected %s Package, Illegal package name \"%s\", at %s", packageType, info.Package, info.Dirname)
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
