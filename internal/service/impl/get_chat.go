package impl

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"context"
	"errors"
)

// GetChat retrieves a chat along with its messages, applying a messages limit.
//
// This method first validates the provided limit string. Then it attempts to fetch
// the chat from the cache. If the chat is not found in cache, it loads the chat
// from storage with the maximum allowed messages, caches it, and then applies
// the requested limit to the messages slice.
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
