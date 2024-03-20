package repository

import "github.com/ellisez/inject-golang/examples/internal/entity"

func EllisAccount() *entity.Account {
	return &entity.Account{
		Username: "ellis",
		Password: "123456",
	}
}
