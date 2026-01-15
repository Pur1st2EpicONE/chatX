package cache

import (
	"chatX/internal/cache/memory"
	"chatX/internal/config"
	"chatX/internal/logger"
	"chatX/internal/models"
)

type Cache interface {
	Get(key int) (models.Chat, error)
	Put(key int, value models.Chat)
	Delete(key int)
	Close()
}

func NewCache(logger logger.Logger, config config.Cache) Cache {
	return memory.NewLRUCache(logger, config)
}
