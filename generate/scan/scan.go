package scan

import (
	"github.com/ellisez/inject-golang/generate/model"
	"github.com/ellisez/inject-golang/generate/parse"
	"os"
	"path/filepath"
)

func DoScan(dirname string) (*model.ModuleInfo, error) {
	p := parse.New()
	p.Result.Dirname = dirname
	err := p.ModParse()
	if err != nil {
		return nil, err
	}
	err = recurDirectory(dirname, p.DoParse)
	if err != nil {
		return nil, err
	}
	return p.Result, nil
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
