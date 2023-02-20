package server

import (
	"context"

	"github.com/rodkevich/mvpbe/pkg/database"
)

// Env represents latent environment configuration for servers in this application.
type Env struct {
	database *database.DB
}

// Option defines function types to modify the Env on creation.
type Option func(*Env) *Env

// NewEnv creates a new Env with the requested options.
func NewEnv(ctx context.Context, opts ...Option) *Env {
	env := &Env{}
	for _, f := range opts {
		env = f(env)
	}

	return env
}

// WithDatabase in the environment.
func WithDatabase(db *database.DB) Option {
	return func(s *Env) *Env {
		s.database = db
		return s
	}
}

// Database ...
func (s *Env) Database() *database.DB {
	return s.database
}

// ShutdownJobs for server env, closing database connections, etc.
func (s *Env) ShutdownJobs(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if s.database != nil {
		s.database.Close(ctx)
	}

	return nil
}
