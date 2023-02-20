package datasource

import (
	"context"
	"fmt"

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

// Readiness of sample database
func (r *SampleDB) Readiness() error {
	return r.db.Pool.Ping(context.Background())
}

// AllDatabases query for all db names
func (r *SampleDB) AllDatabases(ctx context.Context) ([]string, error) {
	sql := `SELECT datname FROM pg_database;`

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
