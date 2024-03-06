package main

import (
	"github.com/ellisez/inject-golang/examples/ctx"
)

// example1
//go:generate inject-golang
// example2 //go:generate inject-golang -m singleton,multiple
// example3 //go:generate inject-golang -m singleton,web github.com/ellisez/inject-golang/examples-work .

func main() {
	c := ctx.New()
	err := c.WebAppStartup1(3001)
	if err != nil {
		return
	}
}
