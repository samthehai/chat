package model

import (
	"time"

	"github.com/samthehai/chat/internal/domain/entity"
)

// Participant model
type Participant struct {
	ID             entity.ID `json:"id"`
	ConversationID entity.ID `json:"conversation_id"`
	UserID         entity.ID `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
