package datasource

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/rodkevich/mvpbe/internal/domain/items_processor/model"
	"github.com/rodkevich/mvpbe/pkg/database"
)

// SampleProcessorDB ...
type SampleProcessorDB struct {
	db *database.DB
}

// New sample model.SampleItem database
func New(db *database.DB) *SampleProcessorDB {
	return &SampleProcessorDB{db: db}
}

// UpdateStatusExampleTrx change item status
func (r *SampleProcessorDB) UpdateStatusExampleTrx(ctx context.Context, m *model.SampleItem) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		sql := `
			UPDATE
			    sample_item
			SET
			    status = $1, end_timestamp=$2
			WHERE
			    item_id = $3;
		`
		resp, err := tx.Exec(ctx, sql, m.Status, m.FinishTime, m.ID)
		if err != nil {
			return fmt.Errorf("UpdateStatusExampleTrx failed: %w", err)
		}
		// must affect at least one row
		if resp.RowsAffected() != 1 {
			return fmt.Errorf("no rows updated")
		}

		return nil
	})
}

// Readiness of sample database
func (r *SampleProcessorDB) Readiness() error {
	return r.db.Pool.Ping(context.Background())
}
