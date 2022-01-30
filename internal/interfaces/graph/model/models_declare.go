package model

import (
	"github.com/samthehai/chat/internal/domain/entity"
)

type CreateNewConversationInput struct {
	Title          string      `json:"title"`
	RecipentIDList []entity.ID `json:"recipentIdList"`
	Text           *string     `json:"text"`
}

type CreateNewConversationPayload struct {
	Conversation *entity.Conversation `json:"conversation"`
}

type PostMessageInput struct {
	ConversationID entity.ID `json:"conversationId"`
	Text           string    `json:"text"`
}

type PostMessagePayload struct {
	Message *entity.Message `json:"message"`
}
