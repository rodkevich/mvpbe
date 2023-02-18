package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/rodkevich/mvpbe/config"
)

func main() {
	cfg := config.NewConfig()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	validator.New()
	r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(cfg)
		if err != nil {
			log.Fatal(fmt.Errorf("JSON encoding failed"))
		}
	})
	log.Fatal(http.ListenAndServe(":"+cfg.API.Port, r))
}
