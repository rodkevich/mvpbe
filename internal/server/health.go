package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rodkevich/mvpbe/pkg/database"
)

// HandleHealth ...
func HandleHealth(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn, err := db.Pool.Acquire(ctx)
		if err != nil {
			log.Println("failed to acquire connection from pool", "error", err)
		}
		defer conn.Release()

		if err := conn.Conn().Ping(ctx); err != nil {
			log.Println("failed to ping database", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok"}`)
	}
}
