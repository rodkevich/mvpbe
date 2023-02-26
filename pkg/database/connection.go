package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/pgxpoolprometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

// DB represents database
type DB struct {
	Pool *pgxpool.Pool
}

// NewPool create new connection using pool
func NewPool(ctx context.Context, connString string, ps *PoolConfig) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	if ps != nil {
		poolCfg.MaxConns = ps.ConnQuantityMax
		poolCfg.MinConns = ps.ConnQuantityMin
		poolCfg.HealthCheckPeriod = time.Duration(ps.HealthCheckPeriod)
		poolCfg.MaxConnIdleTime = time.Duration(ps.ConnTimeIdleMax)
		poolCfg.MaxConnLifetime = time.Duration(ps.ConnTimeLifetime)
	}

	connPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}
	// assert connection works
	err = connPool.Ping(ctx)
	if err != nil {
		log.Fatal("connPool.Ping: ", err.Error())
	}
	// prometheus collector
	collector := pgxpoolprometheus.NewCollector(connPool, map[string]string{"postgres": "postgres"})
	prometheus.MustRegister(collector)

	return &DB{Pool: connPool}, nil
}

// Close pool connections
func (db *DB) Close(_ context.Context) {
	log.Println("Closing db connections pool.")
	db.Pool.Close()
}
