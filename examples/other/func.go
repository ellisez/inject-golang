package other

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/model"
)

func DatabaseLoaded(database *model.Database) {
	fmt.Printf("call Database.postConstruct: %v\n", database)
}
