package sample

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rodkevich/mvpbe/internal/domain/sample/datasource"
	"github.com/rodkevich/mvpbe/internal/middlewares"
	"github.com/rodkevich/mvpbe/internal/server"
	"github.com/rodkevich/mvpbe/internal/serverenv"
)

// Server representation
type Server struct {
	config *Config
	env    *serverenv.ServerEnv
}

// NewServer constructor
func NewServer(cfg *Config, env *serverenv.ServerEnv) (*Server, error) {
	if env.Database() == nil {
		return nil, fmt.Errorf("server requires a database to be presented in the serverenv")
	}

	return &Server{
		env:    env,
		config: cfg,
	}, nil
}

// Routes initialization
func (s *Server) Routes(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.JSONHeaderContentType)

	ds := datasource.New(s.env.Database())
	handler := NewHandler(NewDomain(ds))

	r.Route("/api/v1/sample", func(router chi.Router) {
		router.Get("/health", server.HandleHealth(s.env.Database()))
		router.Get("/liveness", handler.LivenessHandler())
		router.Get("/databases", handler.AllDatabases())
	})

	return r
}
