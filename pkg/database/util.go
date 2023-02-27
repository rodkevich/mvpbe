package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	// register postgres migration driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// register the "file" source migration driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	dbName       = "testing_db-template"
	dbUser       = "testing_db-user"
	dbPassword   = "testing_db-123"
	defaultImage = "postgres:15-alpine"
)

// https://docs.docker.com/language/golang/run-tests/
// https://docs.docker.com/language/golang/configure-ci-cd/

// TestDBInstance is a wrapper around the Docker-based database instance.
type TestDBInstance struct {
	pool       *dockertest.Pool
	container  *dockertest.Resource
	url        *url.URL
	conn       *pgx.Conn
	connLock   sync.Mutex // pg doesn't allow parallel db copy
	skipReason string     // for skip with short or env var setting
}

func MustNewTestInstance() *TestDBInstance {
	testDatabaseInstance, err := NewTestInstance()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return testDatabaseInstance
}

func NewTestInstance() (*TestDBInstance, error) {
	// -short requires flags to be parsed
	if !flag.Parsed() {
		flag.Parse()
	}

	// do not create an instance if -short
	if testing.Short() {
		return &TestDBInstance{
			skipReason: "Skip database tests [-short]!",
		}, nil
	}

	// do not create an instance if SKIP_DB_TESTS is set
	if skip, _ := strconv.ParseBool(os.Getenv("SKIP_DB_TESTS")); skip {
		return &TestDBInstance{
			skipReason: "Skipping database tests [SKIP_DB_TESTS]!",
		}, nil
	}

	ctx := context.Background()

	// create the pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("failed to create db docker pool: %w", err)
	}

	// get the container image to use
	repository, tag, err := postgresRepo()
	if err != nil {
		return nil, fmt.Errorf("failed to determine db repository: %w", err)
	}

	// start the container
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: repository,
		Tag:        tag,
		Env: []string{
			"LANG=C",
			"POSTGRES_DB=" + dbName,
			"POSTGRES_USER=" + dbUser,
			"POSTGRES_PASSWORD=" + dbPassword,
		},
	}, func(c *docker.HostConfig) {
		// set AutoRemove to true
		// so that stopped container goes away by itself
		c.AutoRemove = true
		c.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start db container: %w", err)
	}

	// kill the container in 120 seconds
	if err := container.Expire(120); err != nil {
		return nil, fmt.Errorf("failed to expire db container: %w", err)
	}

	// connection URL
	connectionURL := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(dbUser, dbPassword),
		Host:     container.GetHostPort("5432/tcp"),
		Path:     dbName,
		RawQuery: "sslmode=disable",
	}

	var conn *pgx.Conn
	pool.MaxWait = 5 * time.Minute // retry reads from pg pool MaxWait settings
	if err := pool.Retry(func() error {
		var err error
		conn, err = pgx.Connect(ctx, connectionURL.String())
		if err != nil {
			return err
		}
		if err := conn.Ping(ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed waiting for db container: %w", err)
	}

	// run migrations
	if err := migrateDB(connectionURL.String()); err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	// return the instance
	return &TestDBInstance{
		pool:      pool,
		container: container,
		conn:      conn,
		url:       connectionURL,
	}, nil
}

func (i *TestDBInstance) MustClose() error {
	if err := i.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return nil
}

// Close terminate and clean resources.
func (i *TestDBInstance) Close() (retErr error) {
	// check if we need something to close
	if i.skipReason != "" {
		return
	}

	defer func() {
		if err := i.pool.Purge(i.container); err != nil {
			retErr = fmt.Errorf("failed to purge db container: %w", err)
			return
		}
	}()

	ctx := context.Background()
	if err := i.conn.Close(ctx); err != nil {
		retErr = fmt.Errorf("failed to close connection: %w", err)
		return
	}

	return
}

// NewDatabase creates a new db for tests
func (i *TestDBInstance) NewDatabase(tb testing.TB) (*DB, *Database) {
	tb.Helper() // helper function

	// check if we should create the database
	if i.skipReason != "" {
		tb.Skip(i.skipReason)
	}

	// clone the template database
	newDBName, err := i.clone()
	if err != nil {
		tb.Fatal(err)
	}

	// new connection URL for the new database name
	connectionURL := i.url.ResolveReference(&url.URL{Path: newDBName})
	connectionURL.RawQuery = "sslmode=disable" // resolve will drop this param

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, connectionURL.String())
	if err != nil {
		tb.Fatalf("failed to connect to database %q: %s", newDBName, err)
	}

	// database instance
	db := &DB{Pool: dbpool}

	// close connection and delete db on cleanup
	tb.Cleanup(func() {
		ctx := context.Background()

		// before drop
		db.Close(ctx)

		q := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, newDBName)
		i.connLock.Lock()
		defer i.connLock.Unlock()

		if _, err := i.conn.Exec(ctx, q); err != nil {
			tb.Errorf("failed to drop database %q: %s", newDBName, err)
		}
	})

	host, port, err := net.SplitHostPort(i.url.Host)
	if err != nil {
		tb.Errorf("net.SplitHostPort failed %q: %s", i.url.Host, err)
	}

	return db, &Database{
		Host:     host,
		Port:     port,
		User:     dbUser,
		Password: dbPassword,
		DBName:   newDBName,
		SSLMode:  "disable",
	}
}

// clone creates a new database with a random name from the template instance.
func (i *TestDBInstance) clone() (string, error) {
	// Generate a random database name.
	name, err := randomDatabaseName()
	if err != nil {
		return "", fmt.Errorf("failed to generate random database name: %w", err)
	}

	// usable only for valid prepared statements
	ctx := context.Background()
	q := fmt.Sprintf(`CREATE DATABASE "%s" WITH TEMPLATE "%s";`, name, dbName)

	// postgres does not allow parallel database creation using template
	i.connLock.Lock()
	defer i.connLock.Unlock()

	// clone the template database as the new random database name
	if _, err := i.conn.Exec(ctx, q); err != nil {
		return "", fmt.Errorf("failed to clone template database: %w", err)
	}
	return name, nil
}

func migrateDB(url string) error {
	migrationsDir := fmt.Sprintf("file://%s", migrationsDir())
	m, err := migrate.New(migrationsDir, url)
	if err != nil {
		return fmt.Errorf("failed create migrate: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("migrate source error: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("migrate database error: %w", dbErr)
	}
	return nil
}

func migrationsDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Join(filepath.Dir(filename), "../../migrations")
}

// get from env or default
func postgresRepo() (string, string, error) {
	ref := os.Getenv("CI_POSTGRES_IMAGE")
	if ref == "" {
		ref = defaultImage
	}

	parts := strings.SplitN(ref, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid reference for db container: %q", ref)
	}
	return parts[0], parts[1], nil
}

// randomDatabaseName returns a random database name
func randomDatabaseName() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
