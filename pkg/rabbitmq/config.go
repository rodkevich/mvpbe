package rabbitmq

// Config for amqp [rabbitmq]
type Config struct {
	URI string `default:"amqp://guest:guest@localhost:5672"`
}

// AMQPConfig ...
func (c *Config) AMQPConfig() *Config {
	return c
}
