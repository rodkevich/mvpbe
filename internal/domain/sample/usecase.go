package sample

import (
	"context"
	"errors"
	"fmt"

	"github.com/rodkevich/mvpbe/internal/domain/sample/datasource"
	"github.com/rodkevich/mvpbe/internal/domain/sample/model"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

//go:generate mockery --name ItemsSampleUsage --case underscore  --output mocks/

var errDeletedItem = errors.New("item deleted")

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

func (s *Items) AddItem(ctx context.Context, m *model.SampleItem) error {
	m.StartTime = api.TimeNow()
	m.FinishTime = api.TimeNow()
	m.Status = model.ItemCreated
	return s.itemsRepo.AddItemExampleTrx(ctx, m)
}

func (s *Items) UpdateItem(ctx context.Context, m *model.SampleItem) error {
	m.FinishTime = api.TimeNow()
	return s.itemsRepo.UpdateStatusExampleTrx(ctx, m)
}

func (s *Items) GetItem(ctx context.Context, id string) (*model.SampleItem, error) {
	example, err := s.itemsRepo.GetItemExample(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("remote get failed: %w", err)
	}

	if example.Status == model.ItemDeleted {
		return nil, errDeletedItem
	}
	return example, nil
}

func (s *Items) ListItems(ctx context.Context) ([]*model.SampleItem, error) {
	// TODO implement me
	panic("implement me")
}

// Readiness of domain
func (s *Items) Readiness() error {
	return s.itemsRepo.Readiness()
}

// AllDatabases sample method to get with all db names
func (s *Items) AllDatabases(ctx context.Context) ([]string, error) {
	return s.itemsRepo.AllDatabases(ctx)
}

// NewDomain constructor
func NewDomain(repo *datasource.SampleDB, pbl rabbitmq.AMQPPublisher) *Items {
	return &Items{
		itemsRepo:     repo,
		amqpPublisher: pbl,
	}
}
