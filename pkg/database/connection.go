package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB ...
type DB struct {
	Pool *pgxpool.Pool
}

// New create new connection using pool
func New(ctx context.Context, connString string, ps *PoolConfig) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	poolCfg.MaxConns = ps.ConnQuantityMax
	poolCfg.MinConns = ps.ConnQuantityMin
	poolCfg.HealthCheckPeriod = time.Duration(ps.HealthCheckPeriod)
	poolCfg.MaxConnIdleTime = time.Duration(ps.ConnTimeIdleMax)
	poolCfg.MaxConnLifetime = time.Duration(ps.ConnTimeLifetime)

	connPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}
	return &DB{Pool: connPool}, nil
}

// Close pool connections
func (db *DB) Close(ctx context.Context) {
	log.Println("Closing connection pool.")
	db.Pool.Close()
}
