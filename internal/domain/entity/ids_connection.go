package entity

type IDsConnection struct {
	PageInfo   *PageInfo  `json:"pageInfo"`
	Edges      []*IDsEdge `json:"edges"`
	TotalCount int        `json:"totalCount"`
}

type IDsEdge struct {
	Cursor ID `json:"cursor"`
	Node   ID `json:"node"`
}
