package init

import (
	"github.com/ellisez/inject-golang/examples/model"
	"github.com/gofiber/fiber/v2"
)

// ConfigureWebApp
// @webAppProvide WebApp
// @import github.com/ellisez/inject-golang/examples/model
// @proxy WebAppStartup1
// @injectParam config Config
// @static /images /images
// @static /css /css [Compress,Browse]
// @static /js /js [Compress,Download,Browse] index.html 86400
func ConfigureWebApp(config *model.Config, defaultPort uint) (string, uint, error) {
	if config.Port == 0 {
		defaultPort = config.Port
	}
	return config.Host, defaultPort, nil
}

// WebAppStartupWithHostPort
// @webApp WebApp
// @static /images /images
// @static /css /css [Compress,Browse]
// @static /js /js [Compress,Download,Browse] index.html 86400
func WebAppStartupWithHostPort(webApp *fiber.App, host string, port uint) (string, uint, error) {
	return host, port, nil
}
