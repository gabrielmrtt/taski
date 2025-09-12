package core

type PaginationOutput[T any] struct {
	Data    []T  `json:"data"`
	Page    int  `json:"page"`
	HasMore bool `json:"hasMore"`
	Total   int  `json:"total"`
}

func HasMorePages(currentPage int, totalItems int, perPage int) bool {
	if totalItems <= 0 {
		return false
	}

	return (totalItems / perPage) > currentPage
}
