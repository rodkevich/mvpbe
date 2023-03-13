package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/rodkevich/mvpbe/internal/coverage"
	"github.com/rodkevich/mvpbe/internal/server"
	"github.com/rodkevich/mvpbe/internal/setup"
)

func init() {
	// clean fs from redundant
	cleanFS()

	if err := os.Mkdir(".coverage", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-shuffle=on", "-count=1", "-race", "-timeout=10m", "./...", "-coverprofile=coverage.out")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("could not run go test command: ", err)
	}

	outFile := "coverage.out"
	htmlFile := ".coverage/index.html"
	once := sync.Once{}

	for {
		if fileExists(outFile) {
			once.Do(
				func() {
					fmt.Println("[+] out: file exists")
					convert := exec.Command("go", "tool", "cover", "-html", outFile, "-o", htmlFile)
					convert.Stdout = os.Stdout
					if err := convert.Run(); err != nil {
						fmt.Println("could not run go tool command:", err)
					}
				})
		}
		if fileExists(htmlFile) {
			fmt.Println("[+] index: file exists")
			break
		}
	}

	fmt.Println("Test coverage html generation is finished")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	// check not to be directory while return
	return !info.IsDir()
}

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// recover main goroutine on panic
	defer func() {
		done()
		if r := recover(); r != nil {
			log.Fatal("application got panic", "panic", r)
		}
	}()

	// run app
	err := runCoverageApplication(ctx)
	done()
	if err != nil {
		log.Println(err)
	}

	log.Println("successful shutdown")
}

func runCoverageApplication(ctx context.Context) error {
	// init config
	err := godotenv.Load() // ".coverage.env"
	if err != nil {
		log.Fatal(err)
	}

	var cfg coverage.Config
	envconfig.MustProcess("", &cfg)
	// set up env remotes
	env, err := setup.NewEnvSetup(ctx, &cfg)
	if err != nil {
		return fmt.Errorf("setup.NewEnvSetup: %w", err)
	}

	// set tasks to do on server shut down
	defer func(env *server.Env, ctx context.Context) {
		err := env.ShutdownJobs(ctx)
		if err != nil {
			fmt.Printf("env.ShutdownJobs: %s", err.Error())
		}
		cleanFS()
	}(env, ctx)

	someServer, err := coverage.NewServer(&cfg, env)
	if err != nil {
		return fmt.Errorf("coverage.NewServer: %w", err)
	}

	srv, err := server.New(cfg.HTTP.Port)
	if err != nil {
		return fmt.Errorf("coverage.NewEnv: %w", err)
	}

	log.Println("server listening", "localhost:", cfg.HTTP.Port)
	return srv.ServeHTTPHandler(ctx, someServer.Routes(ctx))
}

func cleanFS() {
	cmd := exec.Command("rm", "coverage.out")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run rm command: ", err)
	}

	cmd = exec.Command("rm", "-rf", ".coverage")
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("could not run rm dir command: ", err)
	}
}
