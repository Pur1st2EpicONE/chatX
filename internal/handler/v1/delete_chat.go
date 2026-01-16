package v1

import "github.com/gin-gonic/gin"

// DeleteChat handles DELETE /chats/:id requests.
//
// Deletes the chat identified by the path parameter ID.
// Responds with statusDeleted on success or an error if deletion fails.
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
