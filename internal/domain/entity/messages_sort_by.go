package entity

type MessagesSortByType string

const (
	MessagesSortByTypeCreatedAt MessagesSortByType = "MESSAGES_SORT_BY_CREATED_AT"
)

func MessagesSortByTypes() []MessagesSortByType {
	return []MessagesSortByType{
		MessagesSortByTypeCreatedAt,
	}
}

func IsValidMessagesSortByType(sortBy string) bool {
	for _, t := range MessagesSortByTypes() {
		if string(t) == sortBy {
			return true
		}
	}

	return false
}
