package impl

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"strconv"
	"strings"
	"unicode/utf8"
)

// validateChat checks whether the provided chat is valid.
//
// It trims whitespace from the chat title, counts its runes, and ensures that
// the title is neither empty nor exceeds the maximum allowed length configured
// in the service.
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

// validateMessage checks whether the provided message is valid.
//
// It trims whitespace from the message text, counts its runes, and ensures that
// the text is neither empty nor exceeds the maximum allowed length configured
// in the service.
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

// validateLimit parses and validates the limit string for retrieving messages.
//
// If the string is empty or "0", it returns the default limit. It ensures the
// limit is a positive integer and does not exceed the configured maximum.
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
