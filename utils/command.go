package utils

import (
	"flag"
	"fmt"
	. "github.com/ellisez/inject-golang/global"
	"os"
	"path/filepath"
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

func CommandParse() error {

	var h bool
	var modeFlag arrayFlags
	var clean bool

	flag.BoolVar(&h, "h", false, "`help`")

	flag.Var(&modeFlag, "m", "Generate `mode`: all (default), singleton, multiple, func, web. example 'singleton,multiple'")

	flag.BoolVar(&clean, "clean", false, "`clean` up all code")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr,
			`inject-golang v0.0.2
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
	if clean {
		err := os.RemoveAll(filepath.Join(GenPackage))
		if err != nil {
			if !os.IsNotExist(err) {
				Failuref("Cleaning failed, %s", err.Error())
			}
		}
		Success("Cleaning completed!")
		os.Exit(0)
	}

	if modeFlag != nil {
		FlagAll = false
		for _, s := range modeFlag {
			switch s {
			case "all":
				FlagAll = true
				break
			case "singleton":
				FlagSingleton = true
				break
			case "multiple":
				FlagMultiple = true
				break
			case "func":
				FlagFunc = true
				break
			case "web":
				FlagWeb = true
				break
			default:
				return fmt.Errorf("Options -m \"%s\" unknown\nUsage: %s", s, flag.Lookup("m").Usage)
			}
		}
	}

	ScanDirectories = flag.Args()

	if len(ScanDirectories) == 0 {
		ScanDirectories = []string{"."}
	}
	return nil
}
