package service

import (
	"github.com/ellisez/inject-golang/examples/internal/entity"
	"github.com/ellisez/inject-golang/examples/internal/repository"
	"github.com/ellisez/inject-golang/examples/internal/vo"
)

// FindAccountByName
// @proxy
// @injectParam database
func FindAccountByName(database *vo.Database, name string) *entity.Account {
	if name == "ellis" {
		return repository.EllisAccount()
	}
	return nil
}
