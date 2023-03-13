package itemsprocessor

import (
	"context"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rodkevich/mvpbe/internal/itemsprocessor/datasource"
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
	log.Println("Configuring routes")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.JSONHeaderContentType)

	r.Mount("/debug", middleware.Profiler())

	ds := datasource.New(s.env.Database())
	itemsUsage := NewItemsDomain(ds, s.env.Publisher())

	runExampleItemsConsumer(ctx, itemsUsage, s.env.Publisher())

	handler := NewItemsHandler(itemsUsage)
	r.Route("/api/v1/items", func(r chi.Router) {
		r.Get("/health", server.HandleHealth(s.env.Database()))
		r.Get("/liveness", handler.LivenessHandler())
		r.Handle("/metrics", promhttp.Handler())
	})
	return r
}

func (s *Server) ConfigureExchange() {
	log.Println("Configuring rabbitmq")
	channel := s.env.Publisher().GetChannel()

	err := channel.ExchangeDeclare(exExchangeName, exExchangeKind, true, false, false, false, nil)
	if err != nil {
		log.Fatal("err := ch.ExchangeDeclare: ", err)
	}

	// configure some ques and their bindings
	for k, v := range map[string]string{
		exQueueNameProcess: exBindingKeyItemsProcessing,
		exQueueNameResults: exBindingKeyItemsReadiness,
	} {
		que, err := channel.QueueDeclare(k, true, false, false, false, nil)
		if err != nil {
			log.Fatal("err := ch.QueueDeclare: ", err)
		}

		err = channel.QueueBind(que.Name, v, exExchangeName, false, nil)
		if err != nil {
			log.Fatal("err := ch.QueueBind: ", err)
		}
	}
}
