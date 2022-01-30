package entity

type ConversationsConnection struct {
	PageInfo   *PageInfo            `json:"pageInfo"`
	Edges      []*ConversationsEdge `json:"edges"`
	TotalCount int                  `json:"totalCount"`
}

type ConversationsEdge struct {
	Cursor ID            `json:"cursor"`
	Node   *Conversation `json:"node"`
}
