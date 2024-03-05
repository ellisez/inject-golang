package main

import (
	"fmt"
	"github.com/ellisez/inject-golang/gen"
	"github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/scan"
	"github.com/ellisez/inject-golang/utils"
	"path/filepath"
)

func init() {
	err := utils.CommandParse()
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
	fmt.Println("at", filepath.Join(global.CurrentDirectory, global.GenPackage))
}
