// Package httpserver provides the concrete implementation of the Server interface
// using Go's standard net/http package with graceful shutdown support.
package httpserver

import (
	"chatX/internal/config"
	"chatX/internal/logger"
	"context"
	"errors"
	"net/http"
	"time"
)

// HttpServer implements the Server interface using net/http.
type HttpServer struct {
	srv             *http.Server  // underlying HTTP server
	shutdownTimeout time.Duration // timeout for graceful shutdown
	logger          logger.Logger // logger instance for structured logging
}

// NewServer creates and configures an HttpServer instance.
func NewServer(logger logger.Logger, config config.Server, handler http.Handler) *HttpServer {

	server := &HttpServer{
		shutdownTimeout: config.ShutdownTimeout,
		logger:          logger,
	}

	server.srv = &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	return server

}

// Run starts the HTTP server and listens for incoming requests.
func (s *HttpServer) Run() error {
	s.logger.LogInfo("server — receiving requests", "layer", "server.httpserver")
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown gracefully stops the HTTP server, waiting up to shutdownTimeout
func (s *HttpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.LogError("server — failed to shutdown gracefully", err, "layer", "server.httpserver")
	} else {
		s.logger.LogInfo("server — shutdown complete", "layer", "server.httpserver")
	}
}
