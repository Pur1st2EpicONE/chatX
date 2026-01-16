// Package models contains the data structures representing
// the domain entities for the chatX application.
package models

import "time"

// Chat represents a chat conversation.
type Chat struct {
	ID        int       `db:"id"`         // Chat ID
	Title     string    `db:"title"`      // Chat title
	CreatedAt time.Time `db:"created_at"` // Chat creation timestamp
	Messages  []Message `db:"messages"`   // Messages in this chat
}

// Message represents a single message in a chat.
type Message struct {
	ID        int       `db:"id"`         // Message ID
	ChatID    int       `db:"chat_id"`    // Parent chat ID
	Text      string    `db:"text"`       // Message text
	CreatedAt time.Time `db:"created_at"` // Message creation timestamp
}
