package model

import (
	"time"

	"github.com/samthehai/chat/internal/domain/entity"
)

// Conversation model
type Conversation struct {
	ID        entity.ID  `json:"id"`
	CreatorID *entity.ID `json:"creator_id"`
	Title     string     `json:"title"`
	Type      string     `json:"type"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func ConvertModelConversation(input *Conversation) *entity.Conversation {
	if input == nil {
		return nil
	}

	return &entity.Conversation{
		ID:        input.ID,
		CreatorID: input.CreatorID,
		Title:     input.Title,
		Type:      entity.ConversationType(input.Type),
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
		DeletedAt: input.DeletedAt,
	}
}

func ConvertModelConversations(input []*Conversation) []*entity.Conversation {
	if input == nil {
		return nil
	}

	cc := make([]*entity.Conversation, 0, len(input))
	for _, c := range input {
		cc = append(cc, ConvertModelConversation(c))
	}

	return cc
}
