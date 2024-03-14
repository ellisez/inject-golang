package model

type Account struct {
	Username string
	Password string
}

func EllisAccount() *Account {
	return &Account{
		Username: "ellis",
		Password: "123456",
	}
}
