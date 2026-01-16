// Package cache provides a generic caching interface
// and factory function for creating cache implementations.
package cache

import (
	"chatX/internal/cache/memory"
	"chatX/internal/config"
	"chatX/internal/logger"
	"chatX/internal/models"
)

// Cache defines the interface for a chat cache.
type Cache interface {
	Get(key int) (models.Chat, error) // Get retrieves a chat by key. Returns ErrCacheMiss if not found.
	Put(key int, value models.Chat)   // Put stores a chat in the cache by key.
	Delete(key int)                   // Delete removes a chat from the cache by key.
	Close()                           // Close releases all cache resources.
}

// NewCache creates a new Cache implementation based on configuration.
func NewCache(logger logger.Logger, config config.Cache) Cache {
	return memory.NewLRUCache(logger, config)
}
