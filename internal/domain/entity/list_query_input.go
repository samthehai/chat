package entity

type ListQueryInput struct {
	First     int
	After     ID
	SortBy    FriendsSortByType
	SortOrder SortOrderType
}
