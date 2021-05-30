package entity

type SortOrderType string

const (
	SortOrderTypeASC SortOrderType = "SORT_ORDER_ASC"
	SortOrderTypeDES SortOrderType = "SORT_ORDER_DES"
)

func sortOrderTypes() []SortOrderType {
	return []SortOrderType{
		SortOrderTypeASC,
		SortOrderTypeDES,
	}
}

func IsValidSortOrderType(soType string) bool {
	for _, t := range sortOrderTypes() {
		if string(t) == soType {
			return true
		}
	}

	return false
}
