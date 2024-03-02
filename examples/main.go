package main

import (
	"github.com/ellisez/inject-golang/examples/ctx"
)

//go:generate inject-golang
func main() {
	c := ctx.New()
	err := c.WebAppStartup()
	if err != nil {
		return
	}

	err = c.WebAppStartupByAddr("", 3000)
	if err != nil {
		return
	}
}
