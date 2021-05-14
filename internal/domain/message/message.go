package message

import (
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	Text      string    `json:"text"`
}
