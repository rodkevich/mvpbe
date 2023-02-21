package setup

import (
	"context"
	"fmt"
	"log"

	"github.com/rodkevich/mvpbe/internal/server"
	"github.com/rodkevich/mvpbe/pkg/configs"
	"github.com/rodkevich/mvpbe/pkg/database"
)

// DatabaseConfigProvider ...
type DatabaseConfigProvider interface{ DatabaseConfig() *database.Config }

// CacheConfigProvider ...
type CacheConfigProvider interface{ CacheConfig() *configs.Cache }

// HTTPConfigProvider ...
type HTTPConfigProvider interface{ HTTPConfig() *configs.HTTP }

// FeaturesConfigProvider ...
type FeaturesConfigProvider interface{ FeaturesConfig() *configs.Features }

// Setup server.Env
func Setup(ctx context.Context, cfg interface{}) (*server.Env, error) {
	var serverEnvOpts []server.Option

	if provider, ok := cfg.(DatabaseConfigProvider); ok {
		log.Println("configuring database")

		_ = provider.DatabaseConfig() // TODO

		dsl := database.ConnectionStringFromEnv()
		p := database.PoolSettingsFromEnv()
		db, err := database.New(ctx, dsl, p)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database: %w", err)
		}

		serverEnvOpts = append(serverEnvOpts, server.WithDatabase(db))
	}
	return server.NewEnv(ctx, serverEnvOpts...), nil
}
