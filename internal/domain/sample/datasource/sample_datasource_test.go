package datasource

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rodkevich/mvpbe/internal/dev"
	"github.com/rodkevich/mvpbe/internal/domain/sample/model"
	"github.com/rodkevich/mvpbe/pkg/database"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

var testDatabaseInstance *database.TestDBInstance

// creates test db for all tests in datasource package
func TestMain(t *testing.M) {
	testDatabaseInstance = database.MustNewTestInstance()
	defer testDatabaseInstance.MustClose()
	t.Run()
}

func TestSampleDB(t *testing.T) {
	t.Parallel()

	tdb, _ := testDatabaseInstance.NewDatabase(t)
	ctx := dev.TestContext(t)

	t.Run("test_database_existence", func(t *testing.T) {
		t.Parallel()

		sampleDB := SampleDB{db: tdb}
		databases, err := sampleDB.AllDatabases(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, databases, "testing_db-template")
	})

	t.Run("should_get_update", func(t *testing.T) {
		t.Parallel()

		sampleDB := SampleDB{db: tdb}
		item := &model.SampleItem{
			StartTime:  api.TimeNow,
			FinishTime: api.TimeNow.Add(5 * time.Minute),
			Status:     model.ItemCreated,
		}

		err := sampleDB.AddItemExampleTrx(ctx, item)
		if err != nil {
			t.Fatal(err)
		}

		item.Status = model.ItemPending
		err = sampleDB.UpdateStatusExampleTrx(ctx, item)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should_not_get_update", func(t *testing.T) {
		t.Parallel()

		sampleDB := SampleDB{db: tdb}
		item := &model.SampleItem{
			StartTime:  api.TimeNow,
			FinishTime: api.TimeNow.Add(15 * time.Minute),
			Status:     model.ItemCreated,
		}

		err := sampleDB.AddItemExampleTrx(ctx, item)
		if err != nil {
			t.Fatal(err)
		}
		item.ID = -1
		err = sampleDB.UpdateStatusExampleTrx(ctx, item)

		assert.Error(t, err, "no rows updated")
	})
}
