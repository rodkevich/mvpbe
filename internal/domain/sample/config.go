package sample

import (
	"github.com/rodkevich/mvpbe/internal/setup"
	"github.com/rodkevich/mvpbe/pkg/configs"
	"github.com/rodkevich/mvpbe/pkg/database"
)

var (
	_ setup.DatabaseConfigProvider = (*Config)(nil)
	_ setup.HTTPConfigProvider     = (*Config)(nil)
	_ setup.CacheConfigProvider    = (*Config)(nil)
	_ setup.FeaturesConfigProvider = (*Config)(nil)
)

// Config for application
type Config struct {
	Database database.Config
	HTTP     configs.HTTP
	Cache    configs.Cache
	Features configs.Features
}

// DatabaseConfig implements setup.DatabaseConfigProvider
func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}

// HTTPConfig implements setup.HTTPConfigProvider
func (c *Config) HTTPConfig() *configs.HTTP {
	return &c.HTTP
}

// CacheConfig implements setup.CacheConfigProvider
func (c *Config) CacheConfig() *configs.Cache {
	return &c.Cache
}

// FeaturesConfig implements setup.FeaturesConfigProvider
func (c *Config) FeaturesConfig() *configs.Features {
	return &c.Features
}
