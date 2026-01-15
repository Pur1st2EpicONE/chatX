package postgres

import (
	"chatX/internal/models"
	"context"
)

func (s *Storage) CreateChat(ctx context.Context, chat *models.Chat) error {
	return s.db.WithContext(ctx).Create(chat).Error
}
