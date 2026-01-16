// Package v1 provides version 1 of the chat API handlers.
package v1

import (
	"chatX/internal/service"
)

const idKey = "id"              // Context key for chat ID
const limitKey = "limit"        // Context key for GET limit
const statusDeleted = "deleted" // Response string for deleted chats

// Handler contains API v1 handlers and holds the service layer.
type Handler struct {
	service service.Service
}

// NewHandler creates a new v1 API handler with the given service layer.
func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}
