package entity

type ConversationMessagesConnection struct {
	PageInfo   *PageInfo                   `json:"pageInfo"`
	Edges      []*ConversationMessagesEdge `json:"edges"`
	TotalCount int                         `json:"totalCount"`
}

type ConversationMessagesEdge struct {
	Cursor ID       `json:"cursor"`
	Node   *Message `json:"node"`
}
