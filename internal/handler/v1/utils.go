package v1

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseChatID extracts and validates the chat ID from the URL path parameter.
//
// Returns the chat ID as an integer, or ErrInvalidChatID if the ID is invalid or non-positive.
func parseChatID(c *gin.Context) (int, error) {
	chatID, err := strconv.Atoi(c.Param(idKey))
	if err != nil || chatID <= 0 {
		return 0, errs.ErrInvalidChatID
	}
	return chatID, nil
}

// mapMessagesToDTO converts a slice of models.Message to a slice of MessageResponseDTO.
//
// Used to format messages for API responses.
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

// respondOK sends a successful HTTP 200 response with a JSON payload.
//
// Wraps the response in a "result" field to maintain consistent API response format.
func respondOK(c *gin.Context, response any) {
	c.JSON(http.StatusOK, gin.H{"result": response})
}

// respondError sends an HTTP error response based on the provided error.
//
// Uses mapErrorToStatus to determine the appropriate status code and message.
func respondError(c *gin.Context, err error) {
	if err != nil {
		status, msg := mapErrorToStatus(err)
		c.AbortWithStatusJSON(status, gin.H{"error": msg})
	}
}

// mapErrorToStatus maps internal application errors to appropriate HTTP status codes.
//
// Returns a tuple of (status code, message) based on the error type.
//   - 400 Bad Request: validation or input errors
//   - 404 Not Found: chat not found
//   - 500 Internal Server Error: all other errors
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
