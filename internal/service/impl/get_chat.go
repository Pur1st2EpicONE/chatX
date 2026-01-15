package impl

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"context"
	"errors"
)

func (s *Service) GetChat(ctx context.Context, chatID int, limitStr string) (models.Chat, error) {

	limit, err := s.validateLimit(limitStr)
	if err != nil {
		return models.Chat{}, err
	}

	chat, err := s.cache.Get(chatID)
	if err != nil {
		chat, err = s.storage.GetChat(ctx, chatID, s.config.GetLimitMax)
		if err != nil {
			if !errors.Is(err, errs.ErrChatNotFound) {
				s.logger.LogError("service — failed to get chat", err, "chatID", chatID, "layer", "service.impl")
			}
			return models.Chat{}, err
		}
		s.cache.Put(chatID, chat)
	}

	if len(chat.Messages) > limit {
		chat.Messages = chat.Messages[:limit]
	}

	return chat, nil

}
