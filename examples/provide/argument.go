package provide

import "fmt"

// ArgInt example for argument instance
// @provide ArgInt argument
// @order "step 1: Setting variable"
func ArgInt() bool {
	fmt.Println("call IsReady.IsReady")
	return false
}
