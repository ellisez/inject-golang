package scan

import (
	. "generate/global"
	"generate/model"
	"generate/parse"
	"os"
	"path/filepath"
)

func DoScan() (*model.AnnotateInfo, error) {
	p := parse.New()
	err := recurDirectory(RootDirectory, p.DoParse)
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
