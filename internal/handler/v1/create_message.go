package v1

import (
	"chatX/internal/errs"
	"chatX/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateMessage handles POST /chats/:id/messages requests.
//
// Expects JSON body with MessageRequestDTO. Returns the created message as MessageResponseDTO.
// Responds with ErrInvalidJSON if JSON parsing fails or error if chat ID is invalid.
func (h *Handler) CreateMessage(c *gin.Context) {

	var dto MessageRequestDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	chatID, err := parseChatID(c)
	if err != nil {
		respondError(c, err)
		return
	}

	msg, err := h.service.CreateMessage(c.Request.Context(), models.Message{ChatID: chatID, Text: dto.Text})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, MessageResponseDTO{
		ID:        msg.ID,
		ChatID:    msg.ChatID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt})

}
