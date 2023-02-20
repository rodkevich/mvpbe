package database

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config ...
type Config struct {
	PoolConfig
	ConnectionConfig
}

// DatabaseConfig implements setup.DatabaseConfigProvider
func (c *Config) DatabaseConfig() *Config {
	return c
}

// PoolConfig presents settings for pgx pool
type PoolConfig struct {
	ConnQuantityMin   int32 `default:"10"`
	ConnQuantityMax   int32 `default:"50"`
	ConnTimeIdleMax   int32 `default:"1"` // minutes
	ConnTimeLifetime  int32 `default:"3"` // minutes
	HealthCheckPeriod int32 `default:"3"` // minutes
}

// ConnectionConfig presents settings for postgres
type ConnectionConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     string `default:"5432"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	DBName   string `default:"postgres"`
	SSLMode  string `default:"disable"` // mode should be either require or disable
}

// ConnectionStringFromEnv for postgres database
func ConnectionStringFromEnv() string {
	s := ConnectionSettingsFromEnv()
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		s.Host, s.Port, s.User, s.DBName, s.SSLMode, s.Password,
	)
}

// ConnectionSettingsFromEnv for postgres database
func ConnectionSettingsFromEnv() *ConnectionConfig {
	var s ConnectionConfig
	err := envconfig.Process("DB", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &s
}

// PoolSettingsFromEnv for postgres database instance
func PoolSettingsFromEnv() *PoolConfig {
	var s PoolConfig
	err := envconfig.Process("DB", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &s
}
