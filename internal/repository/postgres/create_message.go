package postgres

import (
	"chatX/internal/models"
	"context"
)

// CreateMessage inserts a new message record into the database.
func (s *Storage) CreateMessage(ctx context.Context, message *models.Message) error {
	return s.db.WithContext(ctx).Create(message).Error
}
