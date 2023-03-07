package itemsproducer

import (
	"context"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rodkevich/mvpbe/internal/domain/itemsproducer/datasource"
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

	if env.Publisher() == nil {
		return nil, fmt.Errorf("server requires an ampq publisher to be presented in the serverenv")
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
	pbl := s.env.Publisher()

	log.Println("configuring routes")

	items := NewItemsHandler(NewItemsDomain(ds, pbl))
	r.Route("/api/v1/items", func(r chi.Router) {
		r.Get("/health", server.HandleHealth(s.env.Database()))
		r.Get("/liveness", items.LivenessHandler())
		r.Get("/databases", items.AllDatabases())
		r.Handle("/metrics", promhttp.Handler())

		// items:
		r.Post("/", items.CreateItemHandler())
		// r.Get("/", handler.ListItemHandler())
		r.Route("/{itemID}", func(r chi.Router) {
			r.Get("/", items.GetItemHandler())
			r.Put("/", items.UpdateItemHandler())
			// r.Delete("/", handler.DeleteItemHandler())
		})
	})

	return r
}
