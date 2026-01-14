package impl

import (
	"chatX/internal/models"
	"context"
)

func (s *Service) GetChat(ctx context.Context, chatID int, limitStr string) (models.Chat, error) {

	limit, err := s.validateLimit(limitStr)
	if err != nil {
		return models.Chat{}, err
	}

	chat, err := s.storage.GetChat(ctx, chatID, limit)
	if err != nil {
		s.logger.LogError("service — failed to get chat", err, "chatID", chatID, "layer", "service.impl")
		return models.Chat{}, err
	}

	return chat, nil
}
