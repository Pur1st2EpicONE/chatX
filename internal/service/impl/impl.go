// Package impl provides the concrete implementation of the Service interface
// using storage and cache layers.
package impl

import (
	"chatX/internal/cache"
	"chatX/internal/config"
	"chatX/internal/logger"
	"chatX/internal/repository"
)

// Service implements the business logic for managing chats and messages.
type Service struct {
	logger  logger.Logger      // structured logger
	config  config.Service     // service-specific configuration
	cache   cache.Cache        // cache layer for fast access
	storage repository.Storage // persistent storage layer
}

// NewService creates a new Service instance with the provided dependencies.
func NewService(logger logger.Logger, config config.Service, cache cache.Cache, storage repository.Storage) *Service {
	return &Service{logger: logger, cache: cache, config: config, storage: storage}
}
