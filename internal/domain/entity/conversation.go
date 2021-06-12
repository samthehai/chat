package entity

import (
	"time"
)

// Conversation model
type Conversation struct {
	ID        ID               `json:"id"`
	Title     string           `json:"title"`
	CreatorID *ID              `json:"creator_id"`
	Type      ConversationType `json:"type"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt *time.Time       `json:"deleted_at"`
}

type ConversationType string

const (
	ConversationTypeSingle ConversationType = "CONVERSATION_TYPE_SINGLE"
	ConversationTypeGroup  ConversationType = "CONVERSATION_TYPE_GROUP"
)

func conversationTypes() []ConversationType {
	return []ConversationType{
		ConversationTypeSingle,
		ConversationTypeGroup,
	}
}

func IsValidConversationType(ct string) bool {
	for _, t := range conversationTypes() {
		if t == ConversationType(ct) {
			return true
		}
	}

	return false
}
