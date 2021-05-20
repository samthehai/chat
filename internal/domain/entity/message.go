package entity

import (
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	Text      string    `json:"text"`
}
