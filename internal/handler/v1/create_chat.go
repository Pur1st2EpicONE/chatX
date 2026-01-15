package v1

import (
	"chatX/internal/errs"
	"chatX/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateChat(c *gin.Context) {

	var dto ChatRequestDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	chat, err := h.service.CreateChat(c.Request.Context(), models.Chat{Title: dto.Title})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ChatResponseDTO{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt})

}
