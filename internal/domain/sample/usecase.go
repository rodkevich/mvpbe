package sample

import (
	"context"
	"fmt"

	"github.com/rodkevich/mvpbe/internal/domain/sample/datasource"
	"github.com/rodkevich/mvpbe/internal/domain/sample/model"
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
	GetItem(ctx context.Context, id string) (*model.SampleItem, error)
	ListItems(ctx context.Context) ([]*model.SampleItem, error)
}

// Items implements ItemsSampleUsage
type Items struct {
	itemsRepo     *datasource.SampleDB
	amqpPublisher rabbitmq.AMQPPublisher
}

// AddItem ...
func (i *Items) AddItem(ctx context.Context, m *model.SampleItem) error {
	m.StartTime = api.TimeNow()
	m.FinishTime = api.TimeNow()
	m.Status = model.ItemCreated
	return i.itemsRepo.AddItemExampleTrx(ctx, m)
}

// UpdateItem ...
func (i *Items) UpdateItem(ctx context.Context, m *model.SampleItem) error {
	m.FinishTime = api.TimeNow()
	return i.itemsRepo.UpdateStatusExampleTrx(ctx, m)
}

// GetItem ...
func (i *Items) GetItem(ctx context.Context, id string) (*model.SampleItem, error) {
	example, err := i.itemsRepo.GetItemExample(ctx, id)
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
	return i.itemsRepo.Readiness()
}

// AllDatabases sample method to get with all db names
func (i *Items) AllDatabases(ctx context.Context) ([]string, error) {
	return i.itemsRepo.AllDatabases(ctx)
}

// NewDomain constructor
func NewDomain(repo *datasource.SampleDB, pbl rabbitmq.AMQPPublisher) *Items {
	return &Items{
		itemsRepo:     repo,
		amqpPublisher: pbl,
	}
}
