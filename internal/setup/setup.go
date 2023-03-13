package setup

import (
	"context"
	"fmt"
	"log"

	"github.com/rodkevich/mvpbe/internal/server"
	"github.com/rodkevich/mvpbe/pkg/database"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"
	"github.com/rodkevich/mvpbe/pkg/redis"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

// DatabaseConfigProvider ...
type DatabaseConfigProvider interface{ DatabaseConfig() *database.Database }

// CacheConfigProvider ...
type CacheConfigProvider interface{ CacheConfig() *redis.Config }

// HTTPConfigProvider ...
type HTTPConfigProvider interface{ HTTPConfig() *api.Config }

// AMQPConfigProvider ...
type AMQPConfigProvider interface{ AMQPConfig() *rabbitmq.Config }

// NewEnvSetup server.Env
func NewEnvSetup(ctx context.Context, cfg interface{}) (*server.Env, error) {
	var serverEnvOpts []server.Option
	// check if we have to configure db
	if provider, ok := cfg.(DatabaseConfigProvider); ok {
		log.Println("configuring Database")

		conf := provider.DatabaseConfig()
		db, err := database.NewPool(ctx, conf.GetDSN(), conf)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database: %w", err)
		}

		serverEnvOpts = append(serverEnvOpts, server.WithDatabase(db))
	}
	// check if we have to configure amqp
	if provider, ok := cfg.(AMQPConfigProvider); ok {
		log.Println("configuring Amqp")

		conf := provider.AMQPConfig()
		rmq, err := rabbitmq.NewPublisher(conf)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to rabbitmq: %w", err)
		}

		serverEnvOpts = append(serverEnvOpts, server.WithAMQP(rmq))
	}
	// check if we have to configure cache
	if _, ok := cfg.(CacheConfigProvider); ok {
		log.Println("configuring Cache")
	}
	// check if we have to configure http
	if _, ok := cfg.(HTTPConfigProvider); ok {
		log.Println("configuring Http")
	}

	return server.NewEnv(ctx, serverEnvOpts...), nil
}
