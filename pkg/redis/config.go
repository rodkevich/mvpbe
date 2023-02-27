package redis

// Config configuration for ex. redis
type Config struct {
	Host   string `envconfig:"CACHE_HOST" default:"0.0.0.0"`
	Port   string `envconfig:"CACHE_PORT" default:"6379"`
	Name   int    `envconfig:"CACHE_NAME" `
	User   string `envconfig:"CACHE_USER" `
	Pass   string `envconfig:"CACHE_PASS" `
	Time   int    `envconfig:"CACHE_TIME" `
	Enable bool   `envconfig:"CACHE_ENABLE" default:"false"`
}
