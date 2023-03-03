package v1

import (
	"time"
)

// Config configuration
type Config struct {
	Name              string        `envconfig:"HTTP_NAME" default:"mvp_service"`
	Host              string        `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	Port              string        `envconfig:"HTTP_PORT" default:"3080"`
	ReadTimeout       time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"5s"`
	ReadHeaderTimeout time.Duration `envconfig:"HTTP_READ_HEADER_TIMEOUT" default:"5s"`
	WriteTimeout      time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
	IdleTimeout       time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"120s"`
}

// HTTPConfig ...
func (c *Config) HTTPConfig() *Config {
	return c
}
