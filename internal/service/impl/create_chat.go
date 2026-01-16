package impl

import (
	"chatX/internal/models"
	"context"
	"time"
)

// CreateChat creates a new chat in the system.
func (s *Service) CreateChat(ctx context.Context, chat models.Chat) (models.Chat, error) {

	if err := s.validateChat(&chat); err != nil {
		return models.Chat{}, err
	}

	initChat(&chat)

	if err := s.storage.CreateChat(ctx, &chat); err != nil {
		s.logger.LogError("service — failed to create chat", err, "layer", "service.impl")
		return models.Chat{}, err
	}

	return chat, nil

}

// initChat initializes fields for a new chat.
func initChat(chat *models.Chat) {
	chat.CreatedAt = time.Now().UTC()
}
