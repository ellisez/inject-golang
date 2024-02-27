package main

import (
	"github.com/ellisez/inject-golang/examples/ctx"
)

//go:generate -x inject-golang
func main() {
	ctx.New()
}
