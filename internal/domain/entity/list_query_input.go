package entity

type ListQueryInput struct {
	First     int
	After     ID
	SortBy    string
	SortOrder SortOrderType
}
