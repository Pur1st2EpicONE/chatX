// Package server provides abstractions for running the application HTTP server.
package server

import (
	"chatX/internal/config"
	"chatX/internal/logger"
	"chatX/internal/server/httpserver"
	"net/http"
)

// Server defines the interface for running and gracefully shutting down an HTTP server.
type Server interface {
	Run() error // Run starts the server and begins listening for HTTP requests.
	Shutdown()  // Shutdown gracefully stops the server, allowing in-flight requests to complete.
}

// NewServer creates a new Server instance using the internal HTTP server implementation.
func NewServer(logger logger.Logger, config config.Server, handler http.Handler) Server {
	return httpserver.NewServer(logger, config, handler)
}
