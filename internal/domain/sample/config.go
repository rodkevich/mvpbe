package sample

import (
	"github.com/rodkevich/mvpbe/internal/setup"
	"github.com/rodkevich/mvpbe/pkg/configs"
	"github.com/rodkevich/mvpbe/pkg/database"
)

var _ setup.DatabaseConfigProvider = (*Config)(nil)

// DatabaseConfig implements setup.DatabaseConfigProvider
func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}

// Config for application
type Config struct {
	Database database.Config
	HTTP     configs.HTTP
	Cache    configs.Cache
	Features configs.Features
}
