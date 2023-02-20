package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// Server ...
type Server struct {
	ip       string
	port     string
	listener net.Listener
}

// New ...
func New(port string) (*Server, error) {
	addr := fmt.Sprintf(":" + port)
	listener, err := net.Listen("tcp", addr) // to be ready to listen from start
	if err != nil {
		return nil, fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	return &Server{
		ip:       listener.Addr().(*net.TCPAddr).IP.String(),
		port:     strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
		listener: listener,
	}, nil
}

func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()

		log.Println("server.Serve: context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		log.Println("server.Serve: shutting down")
		errCh <- srv.Shutdown(shutdownCtx)
	}()

	// Run the server. This will block until the provided context is closed.
	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	log.Println("server.Serve: serving stopped")

	return nil
}

// ServeHTTPHandler is a convenience wrapper around ServeHTTP
func (s *Server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	return s.ServeHTTP(ctx, &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	})
}

// Addr returns the server's listening address (ip + port).
func (s *Server) Addr() string {
	return net.JoinHostPort(s.ip, s.port)
}

// IP returns the server's listening IP.
func (s *Server) IP() string {
	return s.ip
}

// Port returns the server's listening port.
func (s *Server) Port() string {
	return s.port
}
