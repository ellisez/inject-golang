package router

import "fmt"

// LoginController
// @router /api/login [post]
// @param username query string true 用户名
// @param password query string true 密码
func LoginController(username string, password string) error {
	fmt.Printf("call LoginController: %s, %s\n", username, password)
	return nil
}
