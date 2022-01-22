package model

type FriendsConnection struct{}

func NewFriendsConnection() (*FriendsConnection, error) {
	return &FriendsConnection{}, nil
}

func (fc *FriendsConnection) TotalCount() int {
	// TODO
	return 0
}

func (fc *FriendsConnection) Edges() ([]*FriendsEdge, error) {
	// TODO
	return nil, nil
}

func (fc *FriendsConnection) PageInfo() *PageInfo {
	// TODO
	return nil
}
