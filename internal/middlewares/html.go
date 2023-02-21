package middlewares

import "net/http"

// HTMLHeaderContentType add html header for content type
func HTMLHeaderContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		h.ServeHTTP(w, r)
	})
}
