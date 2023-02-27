package sample

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rodkevich/mvpbe/internal/domain/sample/datasource"
	"github.com/rodkevich/mvpbe/internal/middlewares"
	"github.com/rodkevich/mvpbe/internal/server"
)

// Server representation
type Server struct {
	config *Config
	env    *server.Env
}

// NewServer constructor
func NewServer(cfg *Config, env *server.Env) (*Server, error) {
	if env.Database() == nil {
		return nil, fmt.Errorf("server requires a database to be presented in the serverenv")
	}

	return &Server{
		env:    env,
		config: cfg,
	}, nil
}

// Routes initialization
func (s *Server) Routes(_ context.Context) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.JSONHeaderContentType)

	ds := datasource.New(s.env.Database())
	handler := NewHandler(NewDomain(ds))

	r.Route("/api/v1/sample", func(router chi.Router) {
		router.Get("/health", server.HandleHealth(s.env.Database()))
		router.Get("/liveness", handler.LivenessHandler())

		router.Get("/{id}", handler.GetItemHandler())
		router.Put("/", handler.UpdateItemHandler())
		router.Post("/", handler.CreateItemHandler())

		router.Get("/databases", handler.AllDatabases())
		router.Handle("/metrics", promhttp.Handler())
	})

	return r
}
