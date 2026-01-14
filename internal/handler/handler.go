package handler

import (
	v1 "chatX/internal/handler/v1"
	"chatX/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const templatePath = "web/templates/index.html"

func NewHandler(service service.Service) http.Handler {

	handler := gin.New()

	handler.Use(gin.Recovery())

	apiV1 := handler.Group("/api/v1/chats")
	handlerV1 := v1.NewHandler(service)

	apiV1.POST("/", handlerV1.CreateChat)
	apiV1.POST("/:id/messages/", handlerV1.CreateMessage)

	apiV1.GET("/:id", handlerV1.GetChat)
	apiV1.DELETE("/:id", handlerV1.DeleteChat)

	return handler

}
