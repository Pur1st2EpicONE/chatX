package impl

import (
	"chatX/internal/config"
	"chatX/internal/logger"
	"chatX/internal/repository"
)

type Service struct {
	logger  logger.Logger
	config  config.Service
	storage repository.Storage
}

func NewService(logger logger.Logger, config config.Service, storage repository.Storage) *Service {
	return &Service{logger: logger, config: config, storage: storage}
}
