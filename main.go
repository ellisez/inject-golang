package main

import (
	"fmt"
	"github.com/ellisez/inject-golang/gen"
	. "github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/parse"
	"github.com/ellisez/inject-golang/scan"
	"github.com/ellisez/inject-golang/utils"
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
			utils.Failure("current directory is not a mod, try to run \"go mod init\"")
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
	moduleInfo, err := scan.DoScan()
	if err != nil {
		utils.Failure(err.Error())
	}

	err = gen.DoGen(moduleInfo)
	if err != nil {
		utils.Failure(err.Error())
	}

	utils.Success("Successful!")
	fmt.Println(filepath.Join(Mod.Path, GenPackage), "has generated")
}
