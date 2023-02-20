package setup

import (
	"context"
	"fmt"
	"log"

	"github.com/rodkevich/mvpbe/internal/serverenv"
	"github.com/rodkevich/mvpbe/pkg/database"
)

// DatabaseConfigProvider ...
type DatabaseConfigProvider interface {
	DatabaseConfig() *database.Config
}

// Setup ..
func Setup(ctx context.Context, cfg interface{}) (*serverenv.ServerEnv, error) {
	var serverEnvOpts []serverenv.Option

	if provider, ok := cfg.(DatabaseConfigProvider); ok {
		log.Println("configuring database")

		_ = provider.DatabaseConfig() // TODO
		cs := database.ConnectionStringFromEnv()
		ps := database.PoolSettingsFromEnv()

		db, err := database.New(ctx, cs, ps)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database: %w", err)
		}

		serverEnvOpts = append(serverEnvOpts, serverenv.WithDatabase(db))
	}
	return serverenv.New(ctx, serverEnvOpts...), nil
}
