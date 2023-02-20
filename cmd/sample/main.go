package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/rodkevich/mvpbe/internal/domain/sample"
	"github.com/rodkevich/mvpbe/internal/server"
	"github.com/rodkevich/mvpbe/internal/serverenv"
	"github.com/rodkevich/mvpbe/internal/setup"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// recover main goroutine on panic
	defer func() {
		done()
		if r := recover(); r != nil {
			log.Fatal("application got panic", "panic", r)
		}
	}()

	err := runApplication(ctx)
	done()

	if err != nil {
		log.Println(err)
	}
	log.Println("successful shutdown")
}

func runApplication(ctx context.Context) error {
	// init config
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var cfg sample.Config
	envconfig.MustProcess("HTTP", &cfg.HTTP)
	envconfig.MustProcess("CACHE", &cfg.Cache)
	envconfig.MustProcess("FEATURE", &cfg.Features)

	// set up env remotes
	env, err := setup.Setup(ctx, &cfg)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}

	// set tasks to do on server shut down
	defer func(env *serverenv.ServerEnv, ctx context.Context) {
		err := env.ShutdownJobs(ctx)
		if err != nil {
			fmt.Printf("env.ShutdownJobs: %s", err.Error())
		}
	}(env, ctx)

	someServer, err := sample.NewServer(&cfg, env)
	if err != nil {
		return fmt.Errorf("sample.NewServer: %w", err)
	}

	srv, err := server.New(cfg.HTTP.Port)
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}

	log.Println("server listening", "port", cfg.HTTP.Port)
	return srv.ServeHTTPHandler(ctx, someServer.Routes(ctx))
}
