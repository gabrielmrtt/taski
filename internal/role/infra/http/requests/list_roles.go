package rolehttprequests

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
	roleservice "github.com/gabrielmrtt/taski/internal/role/service"
	"github.com/gabrielmrtt/taski/pkg/stringutils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListRolesRequest struct {
	Name          *string `json:"name" schema:"name"`
	Description   *string `json:"description" schema:"description"`
	Page          *int    `json:"page" schema:"page"`
	PerPage       *int    `json:"perPage" schema:"perPage"`
	SortBy        *string `json:"sortBy" schema:"sortBy"`
	SortDirection *string `json:"sortDirection" schema:"sortDirection"`
	Relations     *string `json:"relations" schema:"relations"`
}

func (r *ListRolesRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListRolesRequest) ToInput() roleservice.ListRolesInput {
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

	var relationsInput core.RelationsInput = make([]string, 0)
	if r.Relations != nil {
		relationsInput = strings.Split(stringutils.CamelCaseToPascalCase(*r.Relations), ",")
	}

	return roleservice.ListRolesInput{
		Filters: rolerepo.RoleFilters{
			Name:        nameFilter,
			Description: descriptionFilter,
		},
		Pagination: core.PaginationInput{
			Page:    r.Page,
			PerPage: r.PerPage,
		},
		SortInput: core.SortInput{
			By:        r.SortBy,
			Direction: &sortDirection,
		},
		RelationsInput: relationsInput,
	}
}
