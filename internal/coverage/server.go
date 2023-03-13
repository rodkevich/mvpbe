package coverage

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
	r.Use(middlewares.HTMLHeaderContentType)

	r.Handle("/*", http.FileServer(http.Dir(".coverage")))
	return r
}
