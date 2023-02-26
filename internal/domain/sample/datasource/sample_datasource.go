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

// New sample database
func New(db *database.DB) *SampleDB {
	return &SampleDB{
		db: db,
	}
}

func (r *SampleDB) InsertExampleTrx(ctx context.Context, it *model.SampleItem) error {
	// panic("not implemented yet")
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		sql := `
			INSERT INTO
				Sample_Item
				(start_timestamp, end_timestamp, status) 
			VALUES 
				($1, $2, $3)
		`
		_, err := tx.Exec(ctx, sql, it.StartTime, it.FinishTime, it.Status)
		if err != nil {
			return fmt.Errorf("InsertExampleTrx failed: %w", err)
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
