package entity

type UserFriendsConnection struct {
	PageInfo   PageInfo
	Edges      []*UserFriendsEdge
	TotalCount int64
}
