package entity

type FriendsConnection struct {
	PageInfo   *PageInfo      `json:"pageInfo"`
	Edges      []*FriendsEdge `json:"edges"`
	TotalCount int            `json:"totalCount"`
}

type FriendsEdge struct {
	Cursor ID    `json:"cursor"`
	Node   *User `json:"node"`
}
