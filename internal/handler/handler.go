// Package handler provides HTTP handlers and routing for the application.
// It sets up API routes, request logging, and Swagger documentation.
package handler

import (
	_ "chatX/docs"
	v1 "chatX/internal/handler/v1"
	"chatX/internal/logger"
	"chatX/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewHandler creates a new HTTP handler with Gin,
// registers API routes for version 1 of the chat API,
// and sets up Swagger documentation.
//
// Parameters:
//   - logger: application logger instance
//   - requestLogging: enable detailed request logging if true
//   - service: business logic service layer
//
// Returns an http.Handler ready to be served.
func NewHandler(logger logger.Logger, requestLogging bool, service service.Service) http.Handler {

	handler := gin.New()
	handler.Use(gin.Recovery())

	if requestLogging {
		handler.Use(middleware(logger))
	}

	apiV1 := handler.Group("/api/v1/chats")
	handlerV1 := v1.NewHandler(service)

	apiV1.POST("/", handlerV1.CreateChat)
	apiV1.POST("/:id/messages/", handlerV1.CreateMessage)

	apiV1.GET("/:id", handlerV1.GetChat)
	apiV1.DELETE("/:id", handlerV1.DeleteChat)

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return handler

}

// middleware returns a Gin middleware that logs requests and response metadata.
//
// Logs include HTTP method, path, query parameters, latency, status code,
// client IP, protocol, user agent, and Gin errors. Logs are categorized as
// Info, Warn, or Error depending on the response status.
func middleware(logger logger.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []any{
			"method", c.Request.Method,
			"path", path,
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
			"status", status,
			"query", query,
			"proto", c.Request.Proto,
			"user_agent", c.Request.UserAgent(),
			"gin_errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"layer", "handler",
		}

		msg := fmt.Sprintf("handler — received %s request to %s", c.Request.Method, path)

		switch status {
		case 500:
			logger.LogError(msg, nil, fields...)
		case 400, 503:
			logger.LogWarn(msg, fields...)
		default:
			logger.LogInfo(msg, fields...)
		}

	}

}
