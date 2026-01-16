package v1

import "time"

// ChatRequestDTO represents the request body for creating a new chat.
type ChatRequestDTO struct {
	Title string `json:"title" example:"The best chat ever!!!"`
}

// ChatResponseDTO represents the response body when a chat is created.
type ChatResponseDTO struct {
	ID        int       `json:"id" example:"1"`
	Title     string    `json:"title" example:"The best chat ever!!!"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-16T12:00:00Z"`
}

// MessageRequestDTO represents the request body for creating a new message.
type MessageRequestDTO struct {
	Text string `json:"text" example:"Hello there!"`
}

// MessageResponseDTO represents the response body for a single message.
type MessageResponseDTO struct {
	ID        int       `json:"id" example:"10"`
	ChatID    int       `json:"chat_id" example:"1"`
	Text      string    `json:"text" example:"Hi!"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-16T12:01:00Z"`
}

// ChatWithMessagesResponseDTO represents a chat along with its messages.
type ChatWithMessagesResponseDTO struct {
	ID        int                  `json:"id" example:"1"`
	Title     string               `json:"title" example:"The best chat ever!!!"`
	CreatedAt time.Time            `json:"created_at" example:"2025-01-16T12:00:00Z"`
	Messages  []MessageResponseDTO `json:"messages"`
}

// OKResponse represents a generic success response with a typed result.
type OKResponse[T any] struct {
	Result T `json:"result"`
}

// InvalidJSONErrorResponse represents a response for invalid JSON input.
type InvalidJSONErrorResponse struct {
	Error string `json:"error" example:"invalid JSON format"`
}

// ValidationErrorResponse represents a response for validation errors.
type ValidationErrorResponse struct {
	Error string `json:"error" example:"message text cannot be empty"`
}

// InvalidChatIDErrorResponse represents a response for invalid chat ID input.
type InvalidChatIDErrorResponse struct {
	Error string `json:"error" example:"invalid chat ID; must be a positive integer"`
}

// InvalidLimitErrorResponse represents a response for invalid limit input.
type InvalidLimitErrorResponse struct {
	Error string `json:"error" example:"invalid limit; must be an integer"`
}

// ChatNotFoundErrorResponse represents a response when a chat is not found.
type ChatNotFoundErrorResponse struct {
	Error string `json:"error" example:"chat not found"`
}

// InternalServerErrorResponse represents a generic internal server error response.
type InternalServerErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}
