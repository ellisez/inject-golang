package main

import (
	"github.com/ellisez/inject-golang/examples/ctx"
)

//go:generate inject-golang
func main() {
	ctx.New()
}
