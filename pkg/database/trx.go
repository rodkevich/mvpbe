package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// InTx runs within a transaction with the isolation level.
func (db *DB) InTx(ctx context.Context, isoLevel pgx.TxIsoLevel, in func(tx pgx.Tx) error) error {
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("get connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel})
	if err != nil {
		return fmt.Errorf("transaction start: %w", err)
	}

	if txError := in(tx); txError != nil {
		if dbError := tx.Rollback(ctx); dbError != nil {
			return fmt.Errorf("transaction rolling back: (db error: %w)", dbError)
		}
		return txError
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit: %w", err)
	}
	return nil
}
