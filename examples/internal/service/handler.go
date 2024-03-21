package service

import (
	"github.com/ellisez/inject-golang/examples/internal/entity"
	"github.com/ellisez/inject-golang/examples/internal/repository"
	"github.com/ellisez/inject-golang/examples/internal/vo"
)

// FindAccountByName
// @proxy
// @import github.com/ellisez/inject-golang/examples/internal/vo
// @import github.com/ellisez/inject-golang/examples/internal/entity en
// @injectParam database
func FindAccountByName(database *vo.Database, name string) *entity.Account {
	if name == "ellis" {
		return repository.EllisAccount()
	}
	return nil
}

// GetEllisAccount
// @proxy
// @import github.com/ellisez/inject-golang/examples/internal/vo
// @import github.com/ellisez/inject-golang/examples/internal/entity en
// @injectParam database
func GetEllisAccount(database *vo.Database) (*entity.Account, error) {
	return repository.EllisAccount(), nil
}
