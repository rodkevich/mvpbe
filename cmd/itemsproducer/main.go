package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/rodkevich/mvpbe/internal/domain/itemsproducer"
	"github.com/rodkevich/mvpbe/internal/server"
	"github.com/rodkevich/mvpbe/internal/setup"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	// recover main goroutine on panic
	defer func() {
		done()
		if r := recover(); r != nil {
			log.Fatal("application got panic", "panic", r)
		}
	}()

	// run app
	err := runSampleApplication(ctx)
	done()
	if err != nil {
		log.Println(err)
	}

	log.Println("successful application shutdown")
}

func runSampleApplication(ctx context.Context) error {
	// init config
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var cfg itemsproducer.Config
	envconfig.MustProcess("", &cfg)

	// set up env remotes
	env, err := setup.NewEnvSetup(ctx, &cfg)
	if err != nil {
		return fmt.Errorf("setup.NewEnvSetup: %w", err)
	}

	// set tasks to do on server shut down
	defer func(ctx context.Context, env *server.Env) {
		err := env.ShutdownJobs(ctx)
		if err != nil {
			fmt.Printf("env.ShutdownJobs: %s", err.Error())
		}
	}(ctx, env)

	itemsProducerServer, err := itemsproducer.NewServer(&cfg, env)
	if err != nil {
		return fmt.Errorf("sample.NewServer: %w", err)
	}

	srv, err := server.New(cfg.HTTP.Port)
	if err != nil {
		return fmt.Errorf("server.NewEnv: %w", err)
	}

	log.Println("server listening", "port", cfg.HTTP.Port)
	return srv.ServeHTTPHandler(ctx, itemsProducerServer.Routes(ctx))
}
