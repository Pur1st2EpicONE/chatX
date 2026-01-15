package postgres

import (
	"chatX/internal/models"
	"context"
)

func (s *Storage) CreateMessage(ctx context.Context, message *models.Message) error {
	return s.db.WithContext(ctx).Create(message).Error
}
