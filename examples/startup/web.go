package startup

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/internal/vo"
)

// ConfigureWebApp
// @webProvide _ _
// @import github.com/ellisez/inject-golang/examples/internal/vo
// @injectParam config Config
// @static /images /images
// @static /css /css [Compress,Browse]
// @static /js /js [Compress,Download,Browse] index.html 86400
func ConfigureWebApp(config *vo.Config, defaultPort uint) (string, error) {
	if config.Port == 0 {
		defaultPort = config.Port
	}
	return fmt.Sprintf("%s:%d", config.Host, defaultPort), nil
}
