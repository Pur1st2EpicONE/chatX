package v1

import "github.com/gin-gonic/gin"

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
