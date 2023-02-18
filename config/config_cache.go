package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Cache configuration
type Cache struct {
	Host string `default:"0.0.0.0"`
	Port string `default:"6379"`
	Name int
	User string
	Pass string
	Time int
}

// CacheConfig processes env to cache configuration
func CacheConfig() Cache {
	var cache Cache
	envconfig.MustProcess("MVP_CACHE", &cache)

	return cache
}
