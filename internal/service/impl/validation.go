package impl

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (s *Service) validateChat(chat *models.Chat) error {

	chat.Title = strings.TrimSpace(chat.Title)
	length := utf8.RuneCountInString(chat.Title)

	if length == 0 {
		return errs.ErrTitleEmpty
	}

	if length > s.config.MaxTitleLength {
		return errs.ErrTitleTooLong
	}

	return nil

}

func (s *Service) validateMessage(message *models.Message) error {

	message.Text = strings.TrimSpace(message.Text)
	length := utf8.RuneCountInString(message.Text)

	if length == 0 {
		return errs.ErrMessageEmpty
	}

	if length > s.config.MaxMessageLength {
		return errs.ErrMessageTooLong
	}

	return nil

}

func (s *Service) validateLimit(limitStr string) (int, error) {

	if limitStr == "" || limitStr == "0" {
		return s.config.GetLimitDefault, nil
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, errs.ErrInvalidLimit
	}

	if limit < 0 {
		return 0, errs.ErrLimitTooSmall
	}

	if limit > s.config.GetLimitMax {
		return 0, errs.ErrLimitTooLarge
	}

	return limit, nil

}
