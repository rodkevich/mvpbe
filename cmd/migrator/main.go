package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/rodkevich/mvpbe/internal/setup"
	"github.com/rodkevich/mvpbe/pkg/database"
)

var (
	sourceFlag       = flag.String("path", "migrations/", "migrations files path")
	migrationTimeout = flag.Duration("timeout", 15*time.Minute, "migration timeout")
)

func main() {
	flag.Parse()

	ctx, done := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer func() {
		done()
		if r := recover(); r != nil {
			log.Fatalf("application got anic: %v", r)
		}
	}()

	err := runMigrations(ctx)
	done()

	if err != nil {
		log.Println(err)
	}
	log.Println("successful shutdown")
}

func runMigrations(ctx context.Context) error {
	// init config
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var config database.Database
	envconfig.MustProcess("", &config)

	env, err := setup.NewEnvSetup(ctx, &config)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	// set tasks to do on server shut down
	defer func() {
		err := env.ShutdownJobs(ctx)
		if err != nil {
			fmt.Printf("env.ShutdownJobs: %s", err.Error())
		}
	}()

	// Run the migrations
	source := fmt.Sprintf("file://%s", *sourceFlag)
	destination, ok := os.LookupEnv("DB_MIGRATE_DESTINATION_URI")
	if !ok {
		log.Println("failed to get migrations url from env: %w", err)
	}

	m, err := migrate.New(source, destination)
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	m.LockTimeout = *migrationTimeout

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed run migrate: %w", err)
	}

	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		return fmt.Errorf("migrate source error: %w", sourceErr)
	}

	if databaseErr != nil {
		return fmt.Errorf("migrate db error: %w", databaseErr)
	}

	log.Println("finished running migrations")

	return nil
}
