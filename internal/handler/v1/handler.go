package v1

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"chatX/internal/service"

	"github.com/gin-gonic/gin"
)

const limitKey = "limit"
const statusDeleted = "deleted"

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

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

func (h *Handler) GetChat(c *gin.Context) {

	chatID, err := parseChatID(c)
	if err != nil {
		respondError(c, err)
		return
	}

	chat, err := h.service.GetChat(c.Request.Context(), chatID, c.Query(limitKey))
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ChatWithMessagesResponseDTO{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
		Messages:  mapMessagesToDTO(chat.Messages)})

}

func (h *Handler) DeleteChat(c *gin.Context) {

	chatID, err := parseChatID(c)
	if err != nil {
		respondError(c, err)
		return
	}

	if err := h.service.DeleteChat(c.Request.Context(), chatID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, statusDeleted)

}
