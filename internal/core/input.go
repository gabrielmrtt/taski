package core

type FileInput struct {
	FileName     string
	FileContent  []byte
	FileMimeType string
}

type SortDirection string

const (
	SortDirectionAsc SortDirection = "asc"
	SortDirectionDec SortDirection = "desc"
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
	GreaterThan *T
	LessThan    *T
	Equals      *T
	Like        *T
	In          *[]T
	Negate      *bool
}

type ServiceInput interface {
	Validate() error
}
