package main

import "github.com/ellisez/inject-golang/examples/ctx/factory"

// example1
//go:generate inject-golang
// example2 //go:generate inject-golang -m singleton,multiple
// example3 //go:generate inject-golang -m singleton,web github.com/ellisez/inject-golang/examples-work .

func main() {
	c := factory.New()
	err := c.WebAppStartup("127.0.0.1", 3001)
	if err != nil {
		return
	}
}
