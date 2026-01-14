package models

import "time"

type Chat struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	Messages  []Message
}

type Message struct {
	ID        int       `db:"id"`
	ChatID    int       `db:"chat_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}
