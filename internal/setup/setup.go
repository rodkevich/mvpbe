package setup

import (
	"context"
	"fmt"
	"log"

	"github.com/rodkevich/mvpbe/internal/server"
	"github.com/rodkevich/mvpbe/pkg/api/v1"
	"github.com/rodkevich/mvpbe/pkg/database"
	"github.com/rodkevich/mvpbe/pkg/features"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"
	"github.com/rodkevich/mvpbe/pkg/redis"
)

// DatabaseConfigProvider ...
type DatabaseConfigProvider interface{ DatabaseConfig() *database.Database }

// CacheConfigProvider ...
type CacheConfigProvider interface{ CacheConfig() *redis.Config }

// HTTPConfigProvider ...
type HTTPConfigProvider interface{ HTTPConfig() *v1.Config }

// FeaturesConfigProvider contains enabled/disabled app features
type FeaturesConfigProvider interface{ FeaturesConfig() *features.Config }

// AMQPConfigProvider ...
type AMQPConfigProvider interface{ AMQPConfig() *rabbitmq.Config }

// NewEnvSetup server.Env
func NewEnvSetup(ctx context.Context, cfg interface{}) (*server.Env, error) {
	var serverEnvOpts []server.Option

	if provider, ok := cfg.(DatabaseConfigProvider); ok {
		log.Println("configuring Database")

		conf := provider.DatabaseConfig()
		db, err := database.NewPool(ctx, conf.GetDSN(), conf)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database: %w", err)
		}

		serverEnvOpts = append(serverEnvOpts, server.WithDatabase(db))
	}
	if _, ok := cfg.(AMQPConfigProvider); ok {
		log.Println("configuring AMQP")
	}
	if _, ok := cfg.(CacheConfigProvider); ok {
		log.Println("configuring Cache")
	}
	if _, ok := cfg.(HTTPConfigProvider); ok {
		log.Println("configuring Http")
	}
	if _, ok := cfg.(FeaturesConfigProvider); ok {
		log.Println("configuring Features")
	}

	return server.NewEnv(ctx, serverEnvOpts...), nil
}
