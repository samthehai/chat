package entity

type PageInfo struct {
	StartCursor     ID
	EndCursor       ID
	HasPreviousPage bool
	HasNextPage     bool
}
