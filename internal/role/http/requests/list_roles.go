package role_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	role_services "github.com/gabrielmrtt/taski/internal/role/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListRolesRequest struct {
	Name          *string `schema:"name"`
	Description   *string `schema:"description"`
	Page          *int    `schema:"page"`
	PerPage       *int    `schema:"per_page"`
	SortBy        *string `schema:"sort_by"`
	SortDirection *string `schema:"sort_direction"`
}

func (r *ListRolesRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListRolesRequest) ToInput() role_services.ListRolesInput {
	var nameFilter *core.ComparableFilter[string] = nil
	if r.Name != nil {
		nameFilter = &core.ComparableFilter[string]{
			Like: r.Name,
		}
	}

	var descriptionFilter *core.ComparableFilter[string] = nil
	if r.Description != nil {
		descriptionFilter = &core.ComparableFilter[string]{
			Like: r.Description,
		}
	}

	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	return role_services.ListRolesInput{
		Filters: role_core.RoleFilters{
			Name:        nameFilter,
			Description: descriptionFilter,
		},
		Pagination: &core.PaginationInput{
			Page:    r.Page,
			PerPage: r.PerPage,
		},
		SortInput: &core.SortInput{
			By:        r.SortBy,
			Direction: &sortDirection,
		},
	}
}
