package postgres

import (
	"chatX/internal/config"
	"chatX/internal/logger"

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

func (s *Storage) GetDB() *gorm.DB {
	return s.db
}
