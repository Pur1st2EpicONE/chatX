// Package repository provides abstractions and implementations for data storage
// operations, including chats and messages. It supports multiple storage backends.
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

// Storage defines the interface for interacting with chat and message data.
type Storage interface {
	CreateChat(ctx context.Context, chat *models.Chat) error                 // CreateChat inserts a new chat into the database.
	CreateMessage(ctx context.Context, message *models.Message) error        // CreateMessage inserts a new message into the database.
	GetChat(ctx context.Context, chatID int, limit int) (models.Chat, error) // GetChat retrieves a chat by ID, optionally limiting the number of messages returned.
	DeleteChat(ctx context.Context, chatID int) error                        // DeleteChat deletes a chat and its messages by chat ID.
	Close()                                                                  // Close closes any resources used by the storage backend (e.g., database connections).
}

// NewStorage creates a new Storage instance using the Postgres backend.
func NewStorage(logger logger.Logger, config config.Storage, db *gorm.DB) Storage {
	return postgres.NewStorage(logger, config, db)
}

// ConnectDB establishes a GORM connection to a PostgreSQL database based on the given configuration.
// Configures connection pool limits and verifies connectivity via Ping.
//
// Returns:
//   - *gorm.DB: a GORM database instance
//   - error: any connection or configuration error
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
