package sample

import (
	"context"

	"github.com/rodkevich/mvpbe/internal/domain/sample/datasource"
)

// UseCase ...
type UseCase interface {
	Readiness() error
	AllDatabases(ctx context.Context) ([]string, error)
}

// Sample implements UseCase
type Sample struct {
	healthRepo *datasource.SampleDB
}

// NewDomain constructor
func NewDomain(repo *datasource.SampleDB) *Sample {
	return &Sample{
		healthRepo: repo,
	}
}

// Readiness of domain
func (u *Sample) Readiness() error {
	return u.healthRepo.Readiness()
}

// AllDatabases sample method to get with all db names
func (u *Sample) AllDatabases(ctx context.Context) ([]string, error) {
	return u.healthRepo.AllDatabases(ctx)
}
