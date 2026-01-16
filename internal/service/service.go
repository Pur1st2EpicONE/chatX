// Package service provides the business logic layer for the chat application,
// defining operations for managing chats and messages.
package service

import (
	"chatX/internal/cache"
	"chatX/internal/config"
	"chatX/internal/logger"
	"chatX/internal/models"
	"chatX/internal/repository"
	"chatX/internal/service/impl"
	"context"
)

// Service defines the interface for chat-related business logic.
type Service interface {
	CreateChat(ctx context.Context, chat models.Chat) (models.Chat, error)             // CreateChat creates a new chat with the given chat data.
	CreateMessage(ctx context.Context, message models.Message) (models.Message, error) // CreateMessage creates a new message in the specified chat.
	GetChat(ctx context.Context, chatID int, limit string) (models.Chat, error)        // GetChat retrieves a chat by ID, optionally limiting the number of messages returned.
	DeleteChat(ctx context.Context, chatID int) error                                  // DeleteChat deletes a chat by ID.
}

// NewService creates a new Service instance using the concrete implementation from the impl package.
func NewService(logger logger.Logger, config config.Service, cache cache.Cache, storage repository.Storage) Service {
	return impl.NewService(logger, config, cache, storage)
}
