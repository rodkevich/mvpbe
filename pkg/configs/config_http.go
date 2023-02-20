package configs

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// HTTP configuration
type HTTP struct {
	Name              string        `default:"mvp_service"`
	Host              string        `default:"0.0.0.0"`
	Port              string        `default:"3080"`
	ReadTimeout       time.Duration `default:"5s"`
	ReadHeaderTimeout time.Duration `default:"5s"`
	WriteTimeout      time.Duration `default:"10s"`
	IdleTimeout       time.Duration `default:"120s"`
	RequestLog        bool          `default:"false"`
}

// HTTPConfig processes env to api configuration
func HTTPConfig() HTTP {
	var api HTTP
	envconfig.MustProcess("API", &api)

	return api
}
