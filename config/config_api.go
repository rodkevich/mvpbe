package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// API configuration
type API struct {
	Name              string        `default:"mvp_service"`
	Host              string        `default:"0.0.0.0"`
	Port              string        `default:"3080"`
	ReadTimeout       time.Duration `default:"5s"`
	ReadHeaderTimeout time.Duration `default:"5s"`
	WriteTimeout      time.Duration `default:"10s"`
	IdleTimeout       time.Duration `default:"120s"`
	RequestLog        bool          `default:"false"`
}

// APIConfig processes env to api configuration
func APIConfig() API {
	var api API
	envconfig.MustProcess("MVP_API", &api)

	return api
}
