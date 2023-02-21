package server

import (
	"log"
	"net/http"

	"github.com/rodkevich/mvpbe/pkg/database"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

// HandleHealth checks
// returns if sever env is ready to process
func HandleHealth(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn, err := db.Pool.Acquire(ctx)
		if err != nil {
			log.Println("failed to get connection from pool", "error", err)
		}
		defer conn.Release()

		if err := conn.Conn().Ping(ctx); err != nil {
			log.Println("failed to ping db", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		api.RenderJSON(w, http.StatusOK, nil)
	}
}
