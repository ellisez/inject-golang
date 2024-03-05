package utils

import (
	"flag"
	"fmt"
	"github.com/ellisez/inject-golang/global"
	"os"
	"strings"
)

type arrayFlags []string

// Value ...
func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}

// Set 方法是flag.Value接口, 设置flag Value的方法.
// 通过多个flag指定的值， 所以我们追加到最终的数组上.
func (i *arrayFlags) Set(value string) error {
	*i = strings.Split(value, ",")
	return nil
}

var h bool
var modeFlag arrayFlags

func CommandParse() error {
	flag.BoolVar(&h, "h", false, "help")

	flag.Var(&modeFlag, "m", "Generate `mode`: all (default), singleton, multiple, func, web. example 'singleton,multiple'")

	modulePath, err := os.Getwd()
	if err != nil {
		return err
	}
	global.CurrentDirectory = modulePath

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr,
			`inject-golang version: 0.0.2
Usage: inject-golang [-h help] [-m mode] pkg1 pkg2 ...
Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if h {
		flag.Usage()
		os.Exit(0)
	}

	if modeFlag != nil {
		for _, s := range modeFlag {
			switch s {
			case "all":
				global.FlagAll = true
				break
			case "singleton":
				global.FlagSingleton = true
				break
			case "multiple":
				global.FlagMultiple = true
				break
			case "func":
				global.FlagFunc = true
				break
			case "web":
				global.FlagWeb = true
				break
			default:
				fmt.Printf("unknown mode \"%s\", %s", s, flag.Lookup("m").Usage)
			}
		}
	}

	global.ScanDirectories = flag.Args()

	if len(global.ScanDirectories) == 0 {
		global.ScanDirectories = []string{modulePath}
	}
	return nil
}
