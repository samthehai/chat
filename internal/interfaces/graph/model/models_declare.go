package model

import (
	"github.com/samthehai/chat/internal/domain/entity"
)

type ConversationMessagesConnection struct {
	PageInfo   *PageInfo                   `json:"pageInfo"`
	Edges      []*ConversationMessagesEdge `json:"edges"`
	TotalCount int                         `json:"totalCount"`
}

type ConversationMessagesEdge struct {
	Cursor entity.ID       `json:"cursor"`
	Node   *entity.Message `json:"node"`
}

type ConversationsConnection struct {
	PageInfo   *PageInfo            `json:"pageInfo"`
	Edges      []*ConversationsEdge `json:"edges"`
	TotalCount int                  `json:"totalCount"`
}

type ConversationsEdge struct {
	Cursor entity.ID            `json:"cursor"`
	Node   *entity.Conversation `json:"node"`
}

type CreateNewConversationInput struct {
	Title          string      `json:"title"`
	RecipentIDList []entity.ID `json:"recipentIdList"`
	Text           *string     `json:"text"`
}

type CreateNewConversationPayload struct {
	Conversation *entity.Conversation `json:"conversation"`
}

type FriendsEdge struct {
	Cursor entity.ID    `json:"cursor"`
	Node   *entity.User `json:"node"`
}

type PostMessageInput struct {
	ConversationID entity.ID `json:"conversationId"`
	Text           string    `json:"text"`
}

type PostMessagePayload struct {
	Message *entity.Message `json:"message"`
}
