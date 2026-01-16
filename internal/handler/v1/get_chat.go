package v1

import "github.com/gin-gonic/gin"

// GetChat handles GET /chats/:id requests.
//
// Retrieves a chat by its ID along with messages, optionally limited by query parameter "limit".
// Responds with ChatWithMessagesResponseDTO on success or an appropriate error if the chat is not found,
// the chat ID is invalid, or other service errors occur.
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
