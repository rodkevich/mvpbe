package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Database configuration
type Database struct {
	Driver            string
	Host              string
	Port              string
	Name              string
	User              string
	Pass              string
	SSLMode           string
	MaxConnectionPool int
}

// DataStoreConfig processes env to database configuration
func DataStoreConfig() Database {
	var db Database
	envconfig.MustProcess("MVP_DB", &db)

	return db
}
