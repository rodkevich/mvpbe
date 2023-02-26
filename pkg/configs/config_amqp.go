package configs

import (
	"github.com/kelseyhightower/envconfig"
)

// Amqp configuration
type Amqp struct {
	URI string `default:"amqp://guest:guest@localhost:5672"`
}

// AmqpConfig processes env to Amqp configuration
func AmqpConfig() Amqp {
	var amqp Amqp
	envconfig.MustProcess("AMQP", &amqp)

	return amqp
}
