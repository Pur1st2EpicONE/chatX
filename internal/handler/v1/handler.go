package v1

import (
	"chatX/internal/service"
)

const idKey = "id"
const limitKey = "limit"
const statusDeleted = "deleted"

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}
