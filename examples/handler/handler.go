package handler

import "github.com/ellisez/inject-golang/examples/model"

// FindAccountByName
// @proxy
// @injectParam database
func FindAccountByName(database *model.Database, name string) *model.Account {
	if name == "ellis" {
		return model.EllisAccount()
	}
	return nil
}
