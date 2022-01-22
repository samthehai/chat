package model

import "github.com/samthehai/chat/internal/domain/entity"

type PageInfo struct {
	StartCursor     entity.ID `json:"startCursor"`
	EndCursor       entity.ID `json:"endCursor"`
	HasPreviousPage bool      `json:"hasPreviousPage"`
	HasNextPage     bool      `json:"hasNextPage"`
}
