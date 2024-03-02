package init

import (
	"github.com/ellisez/inject-golang/examples/model"
	"github.com/gofiber/fiber/v2"
)

// ConfigureWebApp
// @webApp WebApp
// @proxy WebAppStartup1
// @injectParam config Config
// @static /images /images
// @static /css /css [Compress,Browse]
// @static /js /js [Compress,Download,Browse] index.html 86400
func ConfigureWebApp(webApp *fiber.App, config *model.Config) (string, uint, error) {
	return config.Host, config.Port, nil
}

// WebAppStartupWithHostPort
// @webApp WebApp
// @static /images /images
// @static /css /css [Compress,Browse]
// @static /js /js [Compress,Download,Browse] index.html 86400
func WebAppStartupWithHostPort(webApp *fiber.App, host string, port uint) (string, uint, error) {
	return host, port, nil
}
