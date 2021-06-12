package entity

import (
	"time"
)

type Message struct {
	ID             ID          `json:"id"`
	ConversationID ID          `json:"conversation_id"`
	SenderID       ID          `json:"sender_id"`
	Type           MessageType `json:"type"`
	Content        string      `json:"content"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DeletedAt      *time.Time  `json:"deleted_at"`
}

type MessageType string

const (
	MessageTypeText MessageType = "MESSAGE_TYPE_TEXT"
)

func messageTypes() []MessageType {
	return []MessageType{
		MessageTypeText,
	}
}

func IsValidMessageType(mt string) bool {
	for _, t := range messageTypes() {
		if t == MessageType(mt) {
			return true
		}
	}

	return false
}

func (m *Message) GetID() ID {
	return m.ID
}
