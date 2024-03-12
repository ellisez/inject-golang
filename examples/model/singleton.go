package model

import "fmt"

type Config struct {
	Host string
	Port uint
}

type Database struct {
	Host     string
	Port     uint
	Schema   string
	UserName string
	Password string
}

type WebApp struct {
	Config      *Config
	Database    *Database
	MiddleWares []*MiddleWare
	Routers     []*Router
}

// TestLogin example for inject method with uninjected recv
// @proxy
// @injectParam database Database
func (w *WebApp) TestLogin(database *Database) {
	fmt.Printf("call TestLogin: %v, %v\n", w, database)
	for _, router := range w.Routers {
		if router.Path == "/login" {
			err := router.Handle()
			if err != nil {
				return
			}
			break
		}
	}
}
