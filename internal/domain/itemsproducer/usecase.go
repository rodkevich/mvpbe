package itemsproducer

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rodkevich/mvpbe/internal/domain/itemsproducer/datasource"
	"github.com/rodkevich/mvpbe/internal/domain/itemsproducer/model"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

//go:generate mockery --name ItemsSampleUsage --case underscore  --output mocks/

// ItemsSampleUsage represents sample usage of sample domain
type ItemsSampleUsage interface {
	Readiness() error
	AllDatabases(ctx context.Context) ([]string, error)
	AddItem(ctx context.Context, m *model.SampleItem) error
	UpdateItem(ctx context.Context, m *model.SampleItem) error
	GetItem(ctx context.Context, id int) (*model.SampleItem, error)
	ListItems(ctx context.Context) ([]*model.SampleItem, error)
}

// Items implements ItemsSampleUsage
type Items struct {
	db  *datasource.SampleDB
	rmq rabbitmq.AMQPPublisher
}

// AddItem ...
func (i *Items) AddItem(ctx context.Context, m *model.SampleItem) error {
	m.StartTime = api.TimeNow()
	m.FinishTime = api.TimeNow()
	m.Status = model.ItemCreated

	err := i.db.AddItemExampleTrx(ctx, m)
	if err != nil {
		return fmt.Errorf("remote add failed: %w", err)
	}

	dataBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("json.marshal failed: %w", err)
	}

	return i.rmq.Publish(ctx, exExchangeNameItems, exBindingKeyItems,
		amqp.Publishing{
			Headers:   map[string]interface{}{"example-item-trace-id": m.ID},
			Timestamp: api.TimeNow(),
			Body:      dataBytes,
		})
}

// UpdateItem ...
func (i *Items) UpdateItem(ctx context.Context, m *model.SampleItem) error {
	m.FinishTime = api.TimeNow()
	err := i.db.UpdateStatusExampleTrx(ctx, m)
	if err != nil {
		return fmt.Errorf("remote update failed: %w", err)
	}

	dataBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("json.marshal failed: %w", err)
	}

	return i.rmq.Publish(ctx, exExchangeNameItems, exBindingKeyItems,
		amqp.Publishing{
			Headers:   map[string]interface{}{"example-item-trace-id": m.ID},
			Timestamp: api.TimeNow(),
			Body:      dataBytes,
		})
}

// GetItem ...
func (i *Items) GetItem(ctx context.Context, id int) (*model.SampleItem, error) {
	example, err := i.db.GetItemExample(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("remote get failed: %w", err)
	}

	if example.Status == model.ItemDeleted {
		return nil, errDeletedItem
	}
	return example, nil
}

// ListItems ...
func (i *Items) ListItems(ctx context.Context) ([]*model.SampleItem, error) {
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
