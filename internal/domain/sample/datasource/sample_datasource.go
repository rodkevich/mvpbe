package datasource

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/rodkevich/mvpbe/internal/domain/sample/model"
	"github.com/rodkevich/mvpbe/pkg/database"
)

// SampleDB ...
type SampleDB struct {
	db *database.DB
}

// New sample model.SampleItem database
func New(db *database.DB) *SampleDB {
	return &SampleDB{db: db}
}

// InsertExampleTrx ...
func (r *SampleDB) InsertExampleTrx(ctx context.Context, m *model.SampleItem) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		sql := `
			INSERT INTO
				Sample_Item
				(start_timestamp, end_timestamp, status) 
			VALUES 
				($1, $2, $3)
		`
		_, err := tx.Exec(ctx, sql, m.StartTime, m.FinishTime, m.Status)
		if err != nil {
			return fmt.Errorf("InsertExampleTrx failed: %w", err)
		}
		return nil
	})
}

// AddItemExampleTrx ...
func (r *SampleDB) AddItemExampleTrx(ctx context.Context, m *model.SampleItem) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		sql := `
			INSERT INTO
				Sample_Item
				(start_timestamp, end_timestamp, status) 
			VALUES 
				($1, $2, $3)
			RETURNING item_id
		`
		row := tx.QueryRow(ctx, sql, m.StartTime, m.FinishTime, m.Status)
		if err := row.Scan(&m.ID); err != nil {
			return fmt.Errorf("scan item id failed: %w", err)
		}

		return nil
	})
}

// UpdateStatusExampleTrx change item status
func (r *SampleDB) UpdateStatusExampleTrx(ctx context.Context, m *model.SampleItem) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		sql := `
			UPDATE Sample_Item
			SET status = $1
			WHERE item_id = $2
		`
		resp, err := tx.Exec(ctx, sql, m.Status, m.ID)
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
func (r *SampleDB) Readiness() error {
	return r.db.Pool.Ping(context.Background())
}

// AllDatabases query for all db names
func (r *SampleDB) AllDatabases(ctx context.Context) ([]string, error) {
	const sql = `
		SELECT 
		    datname 
		FROM 
		    pg_database;`
	rows, err := r.db.Pool.Query(ctx, sql)
	if err != nil {
		fmt.Printf("Pool.Query: %s", err)
		return nil, err
	}
	defer rows.Close()

	var batch []string
	for rows.Next() {
		var item string
		err = rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		batch = append(batch, item)
	}

	return batch, nil
}
