package postgres

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (s *Storage) GetChat(ctx context.Context, chatID int, limit int) (models.Chat, error) {

	var chat models.Chat

	if err := s.db.WithContext(ctx).Preload("Messages", preload(limit)).First(&chat, chatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Chat{}, errs.ErrChatNotFound
		}
		return models.Chat{}, err
	}

	return chat, nil

}

func preload(limit int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Order("created_at DESC").Limit(limit) }
}
