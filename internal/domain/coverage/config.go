package coverage

import (
	"github.com/rodkevich/mvpbe/pkg/configs"
)

// Config for application
type Config struct {
	HTTP  configs.HTTP
	Cache configs.Cache
}
