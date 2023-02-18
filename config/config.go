package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Config for application
type Config struct {
	API      API
	Database Database
	Cache    Cache
	Features Features
}

// NewConfig process new application config
func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		API:      APIConfig(),
		Database: DataStoreConfig(),
		Cache:    CacheConfig(),
		Features: FeaturesConfig(),
	}
}
