package entity

type FriendsSortByType string

const (
	FriendsSortByTypeName FriendsSortByType = "FRIENDS_SORT_BY_NAME"
)

func friendsSortByTypes() []FriendsSortByType {
	return []FriendsSortByType{
		FriendsSortByTypeName,
	}
}

func IsValidFriendsSortByType(fsType string) bool {
	for _, t := range friendsSortByTypes() {
		if string(t) == fsType {
			return true
		}
	}

	return false
}
