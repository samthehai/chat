package entity

type ConversationsSortByType string

const (
	ConversationsSortByTypeUpdatedAt ConversationsSortByType = "CONVERSATIONS_SORT_BY_UPDATED_AT"
)

func conversationSortByTypes() []ConversationsSortByType {
	return []ConversationsSortByType{
		ConversationsSortByTypeUpdatedAt,
	}
}

func IsValidConversationsSortByType(sortBy string) bool {
	for _, t := range conversationSortByTypes() {
		if string(t) == sortBy {
			return true
		}
	}

	return false
}
