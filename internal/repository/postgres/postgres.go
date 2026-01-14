package postgres

import (
	"chatX/internal/config"
	"chatX/internal/errs"
	"chatX/internal/logger"
	"chatX/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type Storage struct {
	db     *gorm.DB
	logger logger.Logger
	config config.Storage
}

func NewStorage(logger logger.Logger, config config.Storage, db *gorm.DB) *Storage {
	return &Storage{db: db, logger: logger, config: config}
}

func (s *Storage) Close() {
	sqlDB, err := s.db.DB()
	if err != nil {
		s.logger.LogError("postgres — failed to get underlying sql.DB for closing", err, "layer", "repository.postgres")
	} else {
		if err := sqlDB.Close(); err != nil {
			s.logger.LogError("postgres — failed to close properly", err, "layer", "repository.postgres")
		} else {
			s.logger.LogInfo("postgres — database closed", "layer", "repository.postgres")
		}
	}
}

func (s *Storage) CreateChat(ctx context.Context, chat *models.Chat) error {
	return s.db.WithContext(ctx).Create(chat).Error
}

func (s *Storage) CreateMessage(ctx context.Context, message *models.Message) error {
	return s.db.WithContext(ctx).Create(message).Error
}

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
