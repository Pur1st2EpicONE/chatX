package v1

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseChatID(c *gin.Context) (int, error) {
	chatID, err := strconv.Atoi(c.Param(idKey))
	if err != nil || chatID <= 0 {
		return 0, errs.ErrInvalidChatID
	}
	return chatID, nil
}

func mapMessagesToDTO(messages []models.Message) []MessageResponseDTO {

	msgs := make([]MessageResponseDTO, len(messages))

	for i, m := range messages {
		msgs[i] = MessageResponseDTO{
			ID:        m.ID,
			ChatID:    m.ChatID,
			Text:      m.Text,
			CreatedAt: m.CreatedAt,
		}
	}

	return msgs

}

func respondOK(c *gin.Context, response any) {
	c.JSON(http.StatusOK, gin.H{"result": response})
}

func respondError(c *gin.Context, err error) {
	if err != nil {
		status, msg := mapErrorToStatus(err)
		c.AbortWithStatusJSON(status, gin.H{"error": msg})
	}
}

func mapErrorToStatus(err error) (int, string) {

	switch {

	case errors.Is(err, errs.ErrInvalidJSON),
		errors.Is(err, errs.ErrTitleEmpty),
		errors.Is(err, errs.ErrTitleTooLong),
		errors.Is(err, errs.ErrMessageEmpty),
		errors.Is(err, errs.ErrMessageTooLong),
		errors.Is(err, errs.ErrLimitTooSmall),
		errors.Is(err, errs.ErrLimitTooLarge),
		errors.Is(err, errs.ErrInvalidChatID),
		errors.Is(err, errs.ErrInvalidLimit):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, errs.ErrChatNotFound):
		return http.StatusNotFound, err.Error()

	default:
		return http.StatusInternalServerError, errs.ErrInternal.Error()
	}

}
