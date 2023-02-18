package config

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	err := godotenv.Load("testdata/example.env")
	if err != nil {
		t.Fatal(err)
	}
	cfg := &Config{
		API:      APIConfig(),
		Database: DataStoreConfig(),
		Cache:    CacheConfig(),
		Features: FeaturesConfig(),
	}
	assert.Equal(t, "mvp_db", cfg.Database.Name)
}
