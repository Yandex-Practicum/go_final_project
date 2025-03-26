package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultPort            = ":8080"
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 5 * time.Second
)

// HTTPServer represents an HTTP server.
type HTTPServer struct {
	server          *http.Server
	shutdownTimeout time.Duration
	notify          chan error
}

// New creates a new HTTPServer instance with the given handler and options.
func New(handler http.Handler, opts ...Option) *HTTPServer {
	s := &HTTPServer{
		server: &http.Server{
			Handler:      handler,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
		shutdownTimeout: defaultShutdownTimeout,
		notify:          make(chan error),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start starts the HTTP server in a separate goroutine.
func (s *HTTPServer) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Shutdown gracefully shuts down the HTTP server.
func (s *HTTPServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server.Shutdown: %w", err)
	}

	return nil
}

// Notify returns a channel that can be used to receive any errors that occur during server operation.
func (s *HTTPServer) Notify() <-chan error {
	return s.notify
}
