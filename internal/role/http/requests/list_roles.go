package role_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
	role_services "github.com/gabrielmrtt/taski/internal/role/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListRolesRequest struct {
	Name          *string `json:"name" schema:"name"`
	Description   *string `json:"description" schema:"description"`
	Page          *int    `json:"page" schema:"page"`
	PerPage       *int    `json:"per_page" schema:"per_page"`
	SortBy        *string `json:"sort_by" schema:"sort_by"`
	SortDirection *string `json:"sort_direction" schema:"sort_direction"`
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
		Filters: role_repositories.RoleFilters{
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
