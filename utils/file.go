package utils

import (
	"bytes"
	"errors"
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var goModCache string

func ExistsFile(filename string) (bool, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if stat.IsDir() {
		return false, errors.New(fmt.Sprintf("%s is not File!", filename))
	}
	return true, nil
}

func CreateFileIfNotExists(filename string) error {
	exists, err := ExistsFile(filename)
	if err != nil {
		return err
	}
	if !exists {
		_, err = os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExistsDirectory(dirname string) (bool, error) {
	stat, err := os.Stat(dirname)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if !stat.IsDir() {
		return false, errors.New(fmt.Sprintf("%s is not Directory!", dirname))
	}
	return true, nil
}

func CreateDirectoryIfNotExists(dirname string) error {
	exists, err := ExistsDirectory(dirname)
	if err != nil {
		return err
	}
	if !exists {
		err = os.MkdirAll(dirname, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func JoinPath(p string, j ...string) string {
	if !filepath.IsAbs(p) {
		abs, err := filepath.Abs(p)
		if err == nil {
			p = abs
		}
	}
	for _, s := range j {
		if filepath.IsAbs(s) {
			p = s
		} else {
			p = filepath.Join(p, s)
		}
	}
	return p
}

func DirnameOfPackage(packageName string) (string, error) {
	if strings.HasPrefix(packageName, ".") {
		return JoinPath(packageName), nil
	}

	if filepath.IsAbs(packageName) {
		return packageName, nil
	}

	if Mod.Work != nil {
		p := Mod.Work[packageName]
		if p != "" {
			return p, nil
		}
	}

	var version string
	for p, v := range Mod.Require {
		if p == packageName {
			version = v
			break
		}
	}
	if version == "" {
		return "", fmt.Errorf("%s is not found in go.mod, try to \"go get %s\"", packageName, packageName)
	}

	if goModCache == "" {
		out := &bytes.Buffer{}
		command := exec.Command("go", "env", "GOMODCACHE")
		command.Stdout = out
		err := command.Run()
		if err != nil {
			return "", err
		}
		goModCache = out.String()
		goModCache = strings.TrimSuffix(goModCache, "\n")
	}

	return JoinPath(goModCache, packageName+"@"+version), nil
}
