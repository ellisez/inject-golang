package parse

import (
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/utils"
	"go/ast"
	"go/parser"
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
	"strings"
)

func annotateParse(text string) []string {
	prefixLen := len("// ")
	if len(text) < prefixLen {
		return nil
	}
	text = text[prefixLen:]

	var strArr []string
	mode := "seek" // seek | endOfSpace | endOfDoubleQuote | endOfBackQuote
	lastIndex := 0

	add := func(index int) {
		if mode != "seek" {
			strArr = append(strArr, text[lastIndex:index])
			mode = "seek"
		}
		lastIndex = index
	}

	textLen := len(text)
	for index, char := range text {
		switch mode {
		case "seek":
			switch char {
			case ' ', '\t':
				continue
			case '"':
				mode = "endOfDoubleQuote"
				lastIndex = index + 1
				if lastIndex > textLen-1 {
					lastIndex = textLen - 1
				}
			case '`':
				mode = "endOfBackQuote"
				lastIndex = index + 1
				if lastIndex > textLen-1 {
					lastIndex = textLen - 1
				}
			default:
				mode = "endOfSpace"
				lastIndex = index
			}
		case "endOfSpace":
			switch char {
			case ' ', '\t':
				add(index)
			default:
			}
		case "endOfDoubleQuote":
			switch char {
			case '"':
				add(index)
			default:
			}
		case "endOfBackQuote":
			switch char {
			case '`':
				add(index)
			default:
			}
		}
	}
	if mode != "seek" {
		strArr = append(strArr, text[lastIndex:])
	}
	return strArr
}

type Parser struct {
	*model.Module
	Ctx *model.Ctx
}

func ModParse(dirname string) (*model.Module, error) {
	dirname = utils.JoinPath(dirname)
	mod := CacheModMap[dirname]
	if mod != nil {
		return mod, nil
	}

	mod = &model.Module{}

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
	if strings.HasPrefix(dirname, filepath.Join(p.Path, GenPackage)) {
		return nil
	}

	ext := filepath.Ext(filename)
	if ext == ".go" {
		astFile, err := parser.ParseFile(p.Ctx.FileSet, filename, nil, parser.ParseComments)
		if err != nil {
			utils.Failure(err.Error())
		}

		importPath := p.Package
		if p.Path != dirname {
			rel, err := filepath.Rel(p.Path, dirname)
			if err != nil {
				return err
			}
			importPath += "/" + filepath.ToSlash(rel)
		}

		// 解析包信息
		packageName := astFile.Name.String()

		isAllow, packageType := utils.IsAllowedPackageName(importPath, packageName)
		if !isAllow {
			utils.Failuref("Detected %s Package, Illegal package name \"%s\", at %s", packageType, packageName, dirname)
		}

		decls := astFile.Decls

		for _, decl := range decls {
			switch decl.(type) {
			case *ast.FuncDecl:
				funcDecl := decl.(*ast.FuncDecl)
				p.FuncParse(funcDecl, packageName, importPath)

			}
		}

	}
	return nil
}
