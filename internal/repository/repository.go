package repository

import (
	"chatX/internal/config"
	"chatX/internal/logger"
	"chatX/internal/models"
	"chatX/internal/repository/postgres"
	"context"
	"fmt"

	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Storage interface {
	CreateChat(ctx context.Context, chat *models.Chat) error
	CreateMessage(ctx context.Context, message *models.Message) error
	GetChat(ctx context.Context, chatID int, limit int) (models.Chat, error)
	DeleteChat(ctx context.Context, chatID int) error
	Close()
}

func NewStorage(logger logger.Logger, config config.Storage, db *gorm.DB) Storage {
	return postgres.NewStorage(logger, config, db)
}

func ConnectDB(config config.Storage) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode)

	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{Logger: gormLogger.Discard})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB from gorm: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}
