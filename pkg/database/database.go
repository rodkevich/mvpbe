package database

import (
	"context"
	"fmt"

	pgxLib "github.com/jackc/pgx/v5"
)

// InTx runs within a transaction with the isolation level.
func (db *DB) InTx(ctx context.Context, isoLevel pgxLib.TxIsoLevel, in func(tx pgxLib.Tx) error) error {
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("get connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgxLib.TxOptions{IsoLevel: isoLevel})
	if err != nil {
		return fmt.Errorf("transaction start: %w", err)
	}

	if txError := in(tx); txError != nil {
		if rbError := tx.Rollback(ctx); rbError != nil {
			return fmt.Errorf("rolling back transaction: %v (original error: %w)", rbError, txError)
		}
		return txError
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
