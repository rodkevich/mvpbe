package sample

import (
	"context"
	"fmt"
	"log"

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

	if env.Publisher() == nil {
		return nil, fmt.Errorf("server requires an ampq publisher to be presented in the serverenv")
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
	pbl := s.env.Publisher()
	ch := pbl.GetChannel()

	// todo make compact
	log.Println("configuring rabbit ")
	err := ch.ExchangeDeclare(exampleItemsExchangeName, exampleItemsExchangeKind, true, false, false, false, nil)
	if err != nil {
		log.Fatal("err := ch.ExchangeDeclare")
	}
	queue, err := ch.QueueDeclare(exampleItemsQueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("err := ch.QueueDeclare")
	}
	err = ch.QueueBind(queue.Name, exampleItemsBindingKey, exampleItemsExchangeName, false, nil)
	if err != nil {
		log.Fatal("err := ch.QueueBind")
	}

	items := NewItemsHandler(NewItemsDomain(ds, pbl))

	log.Println("configuring routes")

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
