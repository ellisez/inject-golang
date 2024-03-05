package main

import (
	"github.com/ellisez/inject-golang/examples/ctx"
)

//go:generate inject-golang
func main() {
	c := ctx.New()
	err := c.WebAppStartup1(3001)
	if err != nil {
		return
	}
}
