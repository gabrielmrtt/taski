package coredatabase

import (
	"fmt"
	"slices"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/uptrace/bun"
)

// ApplyComparableFilter applies the comparable filter to the bun query
func ApplyComparableFilter[T any](query *bun.SelectQuery, field string, filter *core.ComparableFilter[T]) *bun.SelectQuery {
	if filter.Equals != nil {
		if filter.Negate != nil {
			query.Where(fmt.Sprintf("%s != ?", field), filter.Equals)
		} else {
			query.Where(fmt.Sprintf("%s = ?", field), filter.Equals)
		}
	}

	if filter.Like != nil {
		if filter.Negate != nil {
			query.Where(fmt.Sprintf("%s NOT ILIKE ?", field), filter.Like)
		} else {
			query.Where(fmt.Sprintf("%s ILIKE ?", field), filter.Like)
		}
	}

	if filter.In != nil {
		if filter.Negate != nil {
			query.Where(fmt.Sprintf("%s NOT IN (?)", field), filter.In)
		} else {
			query.Where(fmt.Sprintf("%s IN (?)", field), filter.In)
		}
	}

	if filter.GreaterThan != nil {
		if filter.Negate != nil {
			query.Where(fmt.Sprintf("%s <= ?", field), filter.GreaterThan)
		} else {
			query.Where(fmt.Sprintf("%s > ?", field), filter.GreaterThan)
		}
	}
	if filter.LessThan != nil {
		if filter.Negate != nil {
			query.Where(fmt.Sprintf("%s >= ?", field), filter.LessThan)
		} else {
			query.Where(fmt.Sprintf("%s < ?", field), filter.LessThan)
		}
	}

	return query
}

// ApplyPagination applies the pagination to the bun query
func ApplyPagination(query *bun.SelectQuery, pagination *core.PaginationInput) *bun.SelectQuery {
	var offset int = 0
	var limit int = 10

	if pagination.PerPage != nil {
		limit = *pagination.PerPage
	}

	if pagination.Page != nil {
		offset = (*pagination.Page - 1) * limit
	}

	query.Offset(offset).Limit(limit)

	return query
}

// ApplySort applies the sorts to the bun query in the order of the sorts input
func ApplySort(query *bun.SelectQuery, sorts ...core.SortInput) *bun.SelectQuery {
	for _, sort := range sorts {
		var sortDirection core.SortDirection = core.SortDirectionAsc
		var sortBy *string = sort.By

		if sort.Direction != nil {
			sortDirection = *sort.Direction
		}

		if sortBy != nil {
			query = query.Order(fmt.Sprintf("%s %s", *sortBy, sortDirection))
		}
	}

	return query
}

// ApplyRelations applies the relations to the bun query
func ApplyRelations(query *bun.SelectQuery, relationsInput core.RelationsInput) *bun.SelectQuery {
	var alreadyAppliedRelations []string = make([]string, 0)

	for _, relation := range relationsInput {
		if !slices.Contains(alreadyAppliedRelations, relation) {
			query = query.Relation(relation)
			alreadyAppliedRelations = append(alreadyAppliedRelations, relation)
		}
	}

	return query
}
