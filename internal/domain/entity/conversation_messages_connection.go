package entity

type ConversationMessagesConnection struct {
	PageInfo   PageInfo
	Edges      []*ConversationMessagesEdge
	TotalCount int64
}
