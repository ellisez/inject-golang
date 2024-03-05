package utils

import (
	"flag"
	"fmt"
	"github.com/ellisez/inject-golang/global"
	"os"
	"strings"
)

type arrayValue []string

func (a *arrayValue) String() string {
	return fmt.Sprintf("%v", *a)
}

func (a *arrayValue) Set(s string) error {
	*a = strings.Split(s, ",")
	return nil
}
func CommandParse() error {
	flag.BoolVar(&global.FlagAll, "all", true, "Generate all code")
	flag.BoolVar(&global.FlagSingleton, "singleton", false, "Only Generate singleton code")
	flag.BoolVar(&global.FlagMultiple, "multiple", false, "Only Generate multiple code")
	flag.BoolVar(&global.FlagFunc, "func", false, "Only Generate func code")
	flag.BoolVar(&global.FlagWeb, "web", false, "Only Generate web code")

	modulePath, err := os.Getwd()
	if err != nil {
		return err
	}
	global.CurrentDirectory = modulePath

	var scanDirectories arrayValue
	flag.Var(&scanDirectories, "dirs", "Scan Directories, default \".\"")

	flag.Parse()
	if scanDirectories == nil {
		global.ScanDirectories = []string{modulePath}
	} else {
		global.ScanDirectories = scanDirectories
	}
	return nil
}
