package itemsproducer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rodkevich/mvpbe/internal/dev"
	"github.com/rodkevich/mvpbe/internal/itemsproducer"
	"github.com/rodkevich/mvpbe/internal/itemsproducer/datasource"
	"github.com/rodkevich/mvpbe/internal/itemsproducer/model"
	"github.com/rodkevich/mvpbe/pkg/database"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq/mocks"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

var testDatabaseInstance *database.TestDBInstance

// creates test db for all tests in datasource package
func TestMain(t *testing.M) {
	testDatabaseInstance = database.MustNewTestInstance()
	defer testDatabaseInstance.MustClose()
	t.Run()
}

func TestItems(t *testing.T) {
	t.Parallel()
	rabbitMock := mocks.NewAMQPPublisher(t)
	tdb, _ := testDatabaseInstance.NewDatabase(t)
	ctx := dev.TestContext(t)

	t.Run("should_not_publish", func(t *testing.T) {
		t.Parallel()

		sampleDB := datasource.New(tdb)
		items := itemsproducer.NewItemsDomain(sampleDB, rabbitMock)

		item := &model.SampleItem{
			ManualProc: true,
		}

		err := items.AddOne(ctx, item)
		assert.NoError(t, err)
		rabbitMock.AssertNotCalled(t, "Publish")

		err = items.UpdateOne(ctx, item)
		assert.NoError(t, err)
		rabbitMock.AssertNotCalled(t, "Publish")
	})

	t.Run("should_get_update", func(t *testing.T) {
		t.Parallel()

		sampleDB := datasource.New(tdb)
		items := itemsproducer.NewItemsDomain(sampleDB, rabbitMock)

		item := &model.SampleItem{
			ID:         0,
			StartTime:  api.TimeNow(),
			FinishTime: api.TimeNow().Add(5 * time.Minute),
			Status:     model.ItemCreated,
			ManualProc: true,
		}

		err := items.AddOne(ctx, item)
		assert.NoError(t, err)
		assert.Greater(t, item.ID, 0)

		item.Status = model.ItemPending
		err = items.UpdateOne(ctx, item)
		assert.NoError(t, err)
	})
}
