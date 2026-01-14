package impl

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) CreateMessage(ctx context.Context, message models.Message) (models.Message, error) {

	if err := s.validateMessage(&message); err != nil {
		return models.Message{}, err
	}

	initMessage(&message)

	if err := s.storage.CreateMessage(ctx, &message); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return models.Message{}, errs.ErrChatNotFound
		}
		s.logger.LogError("service — failed to create message", err, "layer", "service.impl")
		return models.Message{}, err
	}

	return message, nil

}

func initMessage(message *models.Message) {
	message.CreatedAt = time.Now().UTC()
}
