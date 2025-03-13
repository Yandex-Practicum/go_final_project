package httpserver

import (
	"net"
	"time"
)

type Option func(*HTTPServer)

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *HTTPServer) {
		s.server.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *HTTPServer) {
		s.server.WriteTimeout = timeout
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *HTTPServer) {
		s.shutdownTimeout = timeout
	}
}

func WithPort(port string) Option {
	return func(s *HTTPServer) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}
