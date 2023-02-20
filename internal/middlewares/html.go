package middlewares

import "net/http"

// HTML ...
func HTML(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		h.ServeHTTP(w, r)
	}
}
