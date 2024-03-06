package scan

import (
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/model"
	"github.com/ellisez/inject-golang/parse"
	"github.com/ellisez/inject-golang/utils"
	"go/token"
	"os"
	"path/filepath"
)

func DoScan() (*model.ModuleInfo, error) {
	moduleInfo := model.NewModuleInfo()
	fileSet := token.NewFileSet()

	for _, directory := range ScanDirectories {
		directory, err := utils.DirnameOfPackage(directory)
		if err != nil {
			return nil, err
		}
		p := &parse.Parser{
			Result:  moduleInfo,
			FileSet: fileSet,
		}

		p.Mod, err = parse.ModParse(directory)
		if err != nil {
			return nil, err
		}
		err = recurDirectory(directory, p.DoParse)
		if err != nil {
			return nil, err
		}
	}
	return moduleInfo, nil
}
func recurDirectory(filename string, handle func(filename string) error) error {
	list, err := os.ReadDir(filename)
	if err != nil {
		return err
	}

	dirs := make([]string, 0)
	for _, fileOrDir := range list {
		subPath := filepath.Join(filename, fileOrDir.Name())
		if fileOrDir.IsDir() {
			dirs = append(dirs, subPath)
		} else {
			err := handle(subPath)
			if err != nil {
				return err
			}
		}
	}

	for _, dirPath := range dirs {
		err := recurDirectory(dirPath, handle)
		if err != nil {
			return err
		}
	}
	return nil
}
