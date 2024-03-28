package examples_work

// RegisterController
// @router /api/register [post]
// @order "2st register router"
// @param username query string true 用户名
// @param gender query string true 性别
// @param age query int true 年龄
func RegisterController(username string, gender string, age int) error {
	return nil
}
