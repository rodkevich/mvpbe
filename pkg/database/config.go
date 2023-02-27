package database

import (
	"fmt"
	"time"
)

// Database contains dsn and pool settings for postgres
type Database struct {
	Driver            string        `envconfig:"DB_DRIVER" default:"postgres"`
	Host              string        `envconfig:"DB_HOST" default:"0.0.0.0"`
	Port              string        `envconfig:"DB_PORT" default:"5432"`
	User              string        `envconfig:"DB_USER" default:"postgres"`
	Password          string        `envconfig:"DB_PASS" default:"postgres"`
	DBName            string        `envconfig:"DB_NAME" default:"postgres"`
	SSLMode           string        `envconfig:"DB_SSLMODE" default:"disable"` // mode should be either require or disable
	ConnQuantityMin   int32         `envconfig:"DB_POOL_MIN_CONNS" default:"10"`
	ConnQuantityMax   int32         `envconfig:"DB_POOL_MAX_CONNS" default:"50"`
	ConnTimeLifetime  time.Duration `envconfig:"DB_POOL_MAX_CONN_LIFETIME" default:"5m"`
	ConnTimeIdleMax   time.Duration `envconfig:"DB_POOL_MAX_CONN_IDLE_TIME" default:"1m"`
	HealthCheckPeriod time.Duration `envconfig:"DB_POOL_HEALTH_CHECK_PERIOD" default:"1m"`
}

// DatabaseConfig implements setup.DatabaseConfigProvider
func (c *Database) DatabaseConfig() *Database {
	return c
}

// GetDSN for postgres database
func (c *Database) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		c.Host, c.Port, c.User, c.DBName, c.SSLMode, c.Password,
	)
}
