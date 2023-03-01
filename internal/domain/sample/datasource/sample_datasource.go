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

// InsertNoReturnExampleTrx doesn't return any info after commit
func (r *SampleDB) InsertNoReturnExampleTrx(ctx context.Context, m *model.SampleItem) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		sql := `
			INSERT INTO
				sample_item
				(start_timestamp, end_timestamp, status)
			VALUES
				($1, $2, $3);
		`
		_, err := tx.Exec(ctx, sql, m.StartTime, m.FinishTime, m.Status)
		if err != nil {
			return fmt.Errorf("InsertNoReturnExampleTrx failed: %w", err)
		}
		return nil
	})
}

// AddItemExampleTrx returns item.ID from db
func (r *SampleDB) AddItemExampleTrx(ctx context.Context, m *model.SampleItem) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		sql := `
			INSERT INTO
				sample_item
				(start_timestamp, end_timestamp, status)
			VALUES
				($1, $2, $3)
			RETURNING item_id;
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
func (r *SampleDB) Readiness() error {
	return r.db.Pool.Ping(context.Background())
}

// AllDatabases query for all db names
func (r *SampleDB) AllDatabases(ctx context.Context) ([]string, error) {
	const sql = `
		SELECT
		    datname
		FROM
		    pg_database;
	`
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

// GetItemExample by id
func (r *SampleDB) GetItemExample(ctx context.Context, id string) (*model.SampleItem, error) {
	sql := `
			SELECT
			    item_id, start_timestamp, end_timestamp, status
			FROM
			    sample_item
			WHERE
			    item_id = $1;
	`
	row := r.db.Pool.QueryRow(ctx, sql, id)

	var m model.SampleItem
	if err := row.Scan(&m.ID, &m.StartTime, &m.FinishTime, &m.Status); err != nil {
		return nil, fmt.Errorf("item id %s: scan failed: %w", id, err)
	}

	return &m, nil
}
