package v1

import "github.com/gin-gonic/gin"

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
