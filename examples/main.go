package main

import "github.com/ellisez/inject-golang/examples/ctx/factory"

// example1 //go:generate inject-golang
// example2 //go:generate inject-golang -m singleton,multiple
// example3
//go:generate inject-golang github.com/ellisez/inject-golang/examples-work .

func main() {
	c := factory.New()
	err := c.WebAppStartup(3001)
	if err != nil {
		return
	}
}
