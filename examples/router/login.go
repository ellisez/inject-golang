package router

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/model"
)

// LoginController
// @router /api/login [post]
// @import github.com/ellisez/inject-golang/examples/model
// @param username query string true 用户名
// @param password query string true 密码
// @injectParam server ServerAlias
func LoginController(username string, password string, server *model.Server) error {
	fmt.Printf("call LoginController: %s, %s\n", username, password)
	server.TriggerEvent("login", map[string]any{
		"username": username,
		"password": password,
	})
	return nil
}
