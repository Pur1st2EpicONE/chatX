package postgres

import (
	"chatX/internal/models"
	"context"
)

// CreateChat inserts a new chat record into the database.
func (s *Storage) CreateChat(ctx context.Context, chat *models.Chat) error {
	return s.db.WithContext(ctx).Create(chat).Error
}
