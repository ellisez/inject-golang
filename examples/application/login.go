package application

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/internal/vo"
)

// LoginController
// @router /api/login [post]
// @order "2st login router"
// @import github.com/ellisez/inject-golang/examples/model
// @param username query string true 用户名
// @param password query string true 密码
// @injectParam server ServerAlias cast
func LoginController(username string, password string, server *vo.Server) error {
	fmt.Printf("call LoginController: %s, %s\n", username, password)
	server.TriggerEvent("login", map[string]any{
		"username": username,
		"password": password,
	})
	return nil
}
