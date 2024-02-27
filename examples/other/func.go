package other

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/model"
)

func DatabaseLoaded(database *model.Database) {
	fmt.Printf("DatabaseLoaded: %v\n", database)
}
