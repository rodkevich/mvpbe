package itemsproducer

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rodkevich/mvpbe/internal/itemsproducer/datasource"
	"github.com/rodkevich/mvpbe/internal/itemsproducer/model"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

//go:generate mockery --name ItemsSampleUsage --case underscore  --output mocks/

// ItemsSampleUsage represents sample usage of sample domain
type ItemsSampleUsage interface {
	Readiness() error
	AllDatabases(ctx context.Context) ([]string, error)
	AddOne(ctx context.Context, m *model.SampleItem) error
	UpdateOne(ctx context.Context, m *model.SampleItem) error
	GetOne(ctx context.Context, id int) (*model.SampleItem, error)
	List(ctx context.Context) ([]*model.SampleItem, error)
}

// Items implements ItemsSampleUsage
type Items struct {
	db  *datasource.SampleDB
	rmq rabbitmq.AMQPPublisher
}

// AddOne ...
func (i *Items) AddOne(ctx context.Context, m *model.SampleItem) error {
	m.StartTime = api.TimeNow()
	m.FinishTime = api.TimeNow()
	m.Status = model.ItemCreated

	err := i.db.AddItemExampleTrx(ctx, m)
	if err != nil {
		return fmt.Errorf("remote add failed: %w", err)
	}

	// return if no need to publish
	// for auto delivery into processing
	if m.ManualProc {
		return nil
	}

	dataBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("json.marshal failed: %w", err)
	}

	return i.rmq.Publish(ctx, exExchangeNameItems, exBindingKeyItems,
		amqp.Publishing{
			Headers: map[string]interface{}{
				"example-item-trace-id": fmt.Sprintf(api.ExampleItemsTracingFormat, m.ID),
			},
			Timestamp: api.TimeNow(),
			Body:      dataBytes,
		})
}

// UpdateOne ...
func (i *Items) UpdateOne(ctx context.Context, m *model.SampleItem) error {
	m.FinishTime = api.TimeNow()

	err := i.db.UpdateStatusExampleTrx(ctx, m)
	if err != nil {
		return fmt.Errorf("remote update failed: %w", err)
	}

	dataBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("json.marshal failed: %w", err)
	}

	// return if no need to publish
	// for auto delivery into processing
	if m.ManualProc {
		return nil
	}

	return i.rmq.Publish(ctx, exExchangeNameItems, exBindingKeyItems,
		amqp.Publishing{
			Headers: map[string]interface{}{
				"example-item-trace-id": fmt.Sprintf(api.ExampleItemsTracingFormat, m.ID),
			},
			Timestamp: api.TimeNow(),
			Body:      dataBytes,
		})
}

// GetOne ...
func (i *Items) GetOne(ctx context.Context, id int) (*model.SampleItem, error) {
	example, err := i.db.GetItemExample(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("remote get failed: %w", err)
	}

	if example.Status == model.ItemDeleted {
		return nil, errDeletedItem
	}
	return example, nil
}

// List ...
func (i *Items) List(ctx context.Context) ([]*model.SampleItem, error) {
	// TODO implement me
	panic("implement me")
}

// Readiness of domain
func (i *Items) Readiness() error {
	return i.db.Readiness()
}

// AllDatabases sample method to get with all db names
func (i *Items) AllDatabases(ctx context.Context) ([]string, error) {
	return i.db.AllDatabases(ctx)
}

// NewItemsDomain constructor
func NewItemsDomain(repo *datasource.SampleDB, pbl rabbitmq.AMQPPublisher) *Items {
	itemsUsage := &Items{db: repo, rmq: pbl}

	return itemsUsage
}
