package provide

import (
	"github.com/ellisez/inject-golang/examples/model"
)

// ConfigureWebApp
// @webProvide
// @import github.com/ellisez/inject-golang/examples/model
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
