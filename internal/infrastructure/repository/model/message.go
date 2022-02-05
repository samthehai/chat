package model

import (
	"time"

	"github.com/samthehai/chat/internal/domain/entity"
)

// Message model
type Message struct {
	ID             entity.ID          `json:"id"`
	ConversationID entity.ID          `json:"conversation_id"`
	SenderID       entity.ID          `json:"sender_id"`
	Type           entity.MessageType `json:"type"`
	Content        string             `json:"content"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DeletedAt      *time.Time         `json:"deleted_at"`
}

func ConvertModelMessage(msg *Message) *entity.Message {
	if msg == nil {
		return nil
	}

	return &entity.Message{
		ID:             msg.ID,
		ConversationID: msg.ConversationID,
		SenderID:       msg.SenderID,
		Type:           msg.Type,
		Content:        msg.Content,
		CreatedAt:      msg.CreatedAt,
		UpdatedAt:      msg.UpdatedAt,
		DeletedAt:      msg.DeletedAt,
	}
}

func ConvertModelMessages(msgs []*Message) []*entity.Message {
	if msgs == nil {
		return nil
	}

	mm := make([]*entity.Message, 0, len(msgs))
	for _, m := range msgs {
		mm = append(mm, ConvertModelMessage(m))
	}

	return mm
}

func GetColumnNameByMessagesSortByType(t entity.MessagesSortByType) string {
	switch t {
	case entity.MessagesSortByTypeCreatedAt:
		return "created_at"
	default:
		return ""
	}
}
