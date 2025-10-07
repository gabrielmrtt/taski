package core

type FileInput struct {
	FileName     string
	FileContent  []byte
	FileMimeType string
}

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

type SortInput struct {
	By        *string
	Direction *SortDirection
}

type PaginationInput struct {
	Page    *int
	PerPage *int
}

type ComparableFilter[T any] struct {
	GreaterThan        *T
	GreaterThanOrEqual *T
	LessThan           *T
	LessThanOrEqual    *T
	Equals             *T
	Like               *T
	In                 *[]T
	Negate             *bool
}

type ServiceInput interface {
	Validate() error
}

type RelationsInput = []string
