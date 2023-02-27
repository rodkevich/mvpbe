package sample

import (
	"context"

	"github.com/rodkevich/mvpbe/internal/domain/sample/datasource"
	"github.com/rodkevich/mvpbe/internal/domain/sample/model"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

//go:generate mockery --name UseCase --case underscore  --output mocks/

// UseCase represents sample usage
type UseCase interface {
	Readiness() error
	AllDatabases(ctx context.Context) ([]string, error)
	CreateItem(ctx context.Context, m *model.SampleItem) error
	UpdateItem(ctx context.Context, m *model.SampleItem) error
	GetItem(ctx context.Context, id string) (*model.SampleItem, error)
	ListItems(ctx context.Context) ([]*model.SampleItem, error)
}

// Sample implements UseCase
type Sample struct {
	healthRepo *datasource.SampleDB
}

func (s *Sample) CreateItem(ctx context.Context, m *model.SampleItem) error {
	m.StartTime = api.TimeNow
	m.FinishTime = api.TimeNow
	m.Status = model.ItemCreated
	return s.healthRepo.AddItemExampleTrx(ctx, m)
}

func (s *Sample) UpdateItem(ctx context.Context, m *model.SampleItem) error {
	m.FinishTime = api.TimeNow
	return s.healthRepo.UpdateStatusExampleTrx(ctx, m)
}

func (s *Sample) GetItem(ctx context.Context, id string) (*model.SampleItem, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Sample) ListItems(ctx context.Context) ([]*model.SampleItem, error) {
	// TODO implement me
	panic("implement me")
}

// Readiness of domain
func (s *Sample) Readiness() error {
	return s.healthRepo.Readiness()
}

// AllDatabases sample method to get with all db names
func (s *Sample) AllDatabases(ctx context.Context) ([]string, error) {
	return s.healthRepo.AllDatabases(ctx)
}

// NewDomain constructor
func NewDomain(repo *datasource.SampleDB) *Sample {
	return &Sample{
		healthRepo: repo,
	}
}
