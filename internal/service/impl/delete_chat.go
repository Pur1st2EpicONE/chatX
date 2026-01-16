package impl

import (
	"chatX/internal/errs"
	"context"
	"errors"
)

// DeleteChat deletes a chat and invalidates its cache entry.
func (s *Service) DeleteChat(ctx context.Context, chatID int) error {
	if err := s.storage.DeleteChat(ctx, chatID); err != nil {
		if !errors.Is(err, errs.ErrChatNotFound) {
			s.logger.LogError("service — failed to delete chat", err, "chatID", chatID, "layer", "service.impl")
		}
		return err
	}
	s.cache.Delete(chatID)
	return nil
}
