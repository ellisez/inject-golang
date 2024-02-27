package main

import (
	"fmt"
	"github.com/ellisez/inject-golang/gen"
	"github.com/ellisez/inject-golang/global"
	"github.com/ellisez/inject-golang/scan"
	"github.com/ellisez/inject-golang/utils"
	"os"
	"path/filepath"
)

func main() {
	modulePath, err := os.Getwd()
	if err != nil {
		utils.Failure(err.Error())
	}

	moduleInfo, err := scan.DoScan(modulePath)
	if err != nil {
		utils.Failure(err.Error())
	}

	err = gen.DoGen(moduleInfo)
	if err != nil {
		utils.Failure(err.Error())
	}

	utils.Success("Successful!")
	fmt.Println("See", filepath.Join(modulePath, global.GenPackage))
}
