package serverenv

import (
	"context"

	"github.com/rodkevich/mvpbe/pkg/database"
)

// ServerEnv represents latent environment configuration for servers in this application.
type ServerEnv struct {
	database *database.DB
}

// Option defines function types to modify the ServerEnv on creation.
type Option func(*ServerEnv) *ServerEnv

// New creates a new ServerEnv with the requested options.
func New(ctx context.Context, opts ...Option) *ServerEnv {
	env := &ServerEnv{}
	for _, f := range opts {
		env = f(env)
	}

	return env
}

// WithDatabase in the environment.
func WithDatabase(db *database.DB) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.database = db
		return s
	}
}

// Database ...
func (s *ServerEnv) Database() *database.DB {
	return s.database
}

// ShutdownJobs for server env, closing database connections, etc.
func (s *ServerEnv) ShutdownJobs(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if s.database != nil {
		s.database.Close(ctx)
	}

	return nil
}
