package entity

type PageInfo struct {
	StartCursor     ID   `json:"startCursor"`
	EndCursor       ID   `json:"endCursor"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	HasNextPage     bool `json:"hasNextPage"`
}
