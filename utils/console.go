package utils

import (
	"fmt"
	"os"
)

var (
	//black        = string([]byte{27, 91, 57, 48, 109})
	red   = string([]byte{27, 91, 57, 49, 109})
	green = string([]byte{27, 91, 57, 50, 109})
	//yellow       = string([]byte{27, 91, 57, 51, 109})
	//blue         = string([]byte{27, 91, 57, 52, 109})
	//magenta      = string([]byte{27, 91, 57, 53, 109})
	//cyan         = string([]byte{27, 91, 57, 54, 109})
	//white        = string([]byte{27, 91, 57, 55, 59, 52, 48, 109})
	reset = string([]byte{27, 91, 48, 109})
	//disableColor = false
)

func Failure(text string) {
	fmt.Println(red + text + reset)
	os.Exit(1)
}

func Failuref(format string, a ...any) {
	Failure(fmt.Sprintf(format, a...))
}

func Success(text string) {
	fmt.Println(green + text + reset)
}

func Successf(format string, a ...any) {
	Success(fmt.Sprintf(format, a...))
}
