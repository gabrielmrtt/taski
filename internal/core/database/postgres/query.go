package core_database_postgres

import (
	"fmt"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/uptrace/bun"
)

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
