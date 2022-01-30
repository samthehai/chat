package entity

import "encoding/json"

type FriendsQueryInput struct {
	UserID ID
	ListQueryInput
}

// String returns a guaranteed unique string that can be used to identify an object
func (fqi FriendsQueryInput) String() string {
	str, err := json.Marshal(fqi)
	if err != nil {
		return ""
	}

	return string(str)
}

// Raw returns the raw, underlaying value of the key
func (fqi FriendsQueryInput) Raw() interface{} {
	return fqi
}
