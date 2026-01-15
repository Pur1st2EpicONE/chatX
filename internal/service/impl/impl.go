package impl

import (
	"chatX/internal/cache"
	"chatX/internal/config"
	"chatX/internal/logger"
	"chatX/internal/repository"
)

type Service struct {
	logger  logger.Logger
	config  config.Service
	cache   cache.Cache
	storage repository.Storage
}

func NewService(logger logger.Logger, config config.Service, cache cache.Cache, storage repository.Storage) *Service {
	return &Service{logger: logger, cache: cache, config: config, storage: storage}
}
