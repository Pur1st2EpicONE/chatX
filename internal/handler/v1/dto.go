package v1

import "time"

type ChatRequestDTO struct {
	Title string `json:"title" validate:"required"`
}

type ChatResponseDTO struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageRequestDTO struct {
	Text string `json:"text" validate:"required"`
}

type MessageResponseDTO struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatWithMessagesResponseDTO struct {
	ID        int                  `json:"id"`
	Title     string               `json:"title"`
	CreatedAt time.Time            `json:"created_at"`
	Messages  []MessageResponseDTO `json:"messages"`
}
