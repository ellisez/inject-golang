package main

import (
	"github.com/ellisez/inject-golang/generate/gen"
	"github.com/ellisez/inject-golang/generate/scan"
	"os"
)

func main() {
	modulePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	moduleInfo, err := scan.DoScan(modulePath)
	if err != nil {
		panic(err)
	}

	err = gen.DoGen(moduleInfo)
	if err != nil {
		panic(err)
	}

}
