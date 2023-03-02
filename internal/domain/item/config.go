package item

import (
	"github.com/rodkevich/mvpbe/internal/setup"
	"github.com/rodkevich/mvpbe/pkg/database"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"
	"github.com/rodkevich/mvpbe/pkg/redis"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

var (
	_ setup.DatabaseConfigProvider = (*Config)(nil)
	_ setup.HTTPConfigProvider     = (*Config)(nil)
	_ setup.CacheConfigProvider    = (*Config)(nil)
	_ setup.AMQPConfigProvider     = (*Config)(nil)
)

const (
	// rabbit exchange settings todo move to cfg
	exampleItemsQueueName       = "example_items"
	exampleItemsExchangeName    = "example_items_exchange"
	exampleItemsBindingKey      = "example_items_binding_key"
	exampleItemsExchangeKind    = "direct"
	exampleItemsAMQPConcurrency = 10
)

// Config for application
type Config struct {
	AMQP     rabbitmq.Config
	Cache    redis.Config
	Database database.Database
	HTTP     api.Config
}

// DatabaseConfig implements setup.DatabaseConfigProvider
func (c *Config) DatabaseConfig() *database.Database {
	return &c.Database
}

// AMQPConfig implements setup.AMQPConfigProvider
func (c *Config) AMQPConfig() *rabbitmq.Config {
	return &c.AMQP
}

// HTTPConfig implements setup.HTTPConfigProvider
func (c *Config) HTTPConfig() *api.Config {
	return &c.HTTP
}

// CacheConfig implements setup.CacheConfigProvider
func (c *Config) CacheConfig() *redis.Config {
	return &c.Cache
}
