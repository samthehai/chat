package entity

import (
	"time"
)

// Participant model
type Participant struct {
	ID             ID        `json:"id"`
	ConversationID ID        `json:"conversation_id"`
	UserID         ID        `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
