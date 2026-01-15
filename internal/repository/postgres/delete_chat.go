package postgres

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"context"
)

func (s *Storage) DeleteChat(ctx context.Context, chatID int) error {

	result := s.db.WithContext(ctx).Delete(&models.Chat{}, chatID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrChatNotFound
	}

	return nil

}
