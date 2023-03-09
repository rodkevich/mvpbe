package itemsprocessor

import (
	"github.com/rodkevich/mvpbe/internal/setup"
	"github.com/rodkevich/mvpbe/pkg/database"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

var (
	_ setup.DatabaseConfigProvider = (*Config)(nil)
	_ setup.HTTPConfigProvider     = (*Config)(nil)
	_ setup.AMQPConfigProvider     = (*Config)(nil)
)

const (
	// example rabbit settings // TODO move to cfg
	exQueueNameItems           = "example_items"
	exQueueNameResults         = "example_results"
	exExchangeNameItems        = "example_items_exchange"
	exBindingKeyItems          = "example_items_binding_key"
	exBindingKeyItemsReadiness = "example_items_binding_readiness_key"
	exConsumerNameItems        = "items_processor"
	exExchangeKindItems        = "direct"
	exAMQPConcurrencyItems     = 10
)

// Config for application
type Config struct {
	AMQP     rabbitmq.Config
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
