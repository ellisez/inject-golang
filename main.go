package main

import (
	"errors"
	"fmt"
	"github.com/ellisez/inject-golang/gen"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/parse"
	"github.com/ellisez/inject-golang/scan"
	"github.com/ellisez/inject-golang/utils"
	"io/fs"
	"os"
	"path/filepath"
)

func init() {
	modulePath, err := os.Getwd()
	if err != nil {
		utils.Failure(err.Error())
	}

	mod, err := parse.ModParse(modulePath)
	if err != nil {
		if os.IsNotExist(err) {
			var pathError *fs.PathError
			errors.As(err, &pathError)
			filename := pathError.Path
			utils.Failuref(`"%s" is not a mod, try to run "go mod init"`, filepath.Dir(filename))
		} else {
			utils.Failure(err.Error())
		}
	}
	Mod = mod

	err = utils.CommandParse()
	if err != nil {
		utils.Failure(err.Error())
	}
}

func main() {
	ctx, err := scan.DoScan()
	if err != nil {
		utils.Failure(err.Error())
	}

	err = gen.DoGen(ctx)
	if err != nil {
		utils.Failure(err.Error())
	}

	utils.Success("Successful!")
	fmt.Println(filepath.Join(Mod.Path, GenPackage), "generated")
}
