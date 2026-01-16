// Package postgres provides a PostgreSQL implementation of the Storage interface.
package postgres

import (
	"chatX/internal/config"
	"chatX/internal/logger"

	"gorm.io/gorm"
)

// Storage implements the repository.Storage interface for PostgreSQL using GORM.
type Storage struct {
	db     *gorm.DB       // underlying GORM DB connection
	logger logger.Logger  // logger instance for structured logging
	config config.Storage // configuration for database connection
}

// NewStorage creates a new Postgres storage instance.
func NewStorage(logger logger.Logger, config config.Storage, db *gorm.DB) *Storage {
	return &Storage{db: db, logger: logger, config: config}
}

// Close closes the underlying SQL database connection.
// Logs errors if the connection cannot be closed properly.
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
