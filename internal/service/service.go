package service

import (
	"chatX/internal/config"
	"chatX/internal/logger"
	"chatX/internal/models"
	"chatX/internal/repository"
	"chatX/internal/service/impl"
	"context"
)

type Service interface {
	CreateChat(ctx context.Context, chat models.Chat) (models.Chat, error)
	CreateMessage(ctx context.Context, message models.Message) (models.Message, error)
	GetChat(ctx context.Context, chatID int, limit string) (models.Chat, error)
	DeleteChat(ctx context.Context, chatID int) error
}

func NewService(logger logger.Logger, config config.Service, storage repository.Storage) Service {
	return impl.NewService(logger, config, storage)
}
