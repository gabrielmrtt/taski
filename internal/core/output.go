package core

type PaginationOutput[T any] struct {
	Data    []T
	Page    int
	HasMore bool
	Total   int
}

func HasMorePages(currentPage int, totalItems int, perPage int) bool {
	if totalItems <= 0 {
		return false
	}

	return (totalItems / perPage) > currentPage
}
