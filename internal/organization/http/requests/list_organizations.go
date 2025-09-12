package organization_http_requests

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	"github.com/gabrielmrtt/taski/pkg/stringutils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListOrganizationsRequest struct {
	Name          *string `schema:"name"`
	Status        *string `schema:"status"`
	Page          *int    `schema:"page"`
	PerPage       *int    `schema:"perPage"`
	SortBy        *string `schema:"sortBy"`
	SortDirection *string `schema:"sortDirection"`
	Relations     *string `schema:"relations"`
}

func (r *ListOrganizationsRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListOrganizationsRequest) ToInput() organization_services.ListOrganizationsInput {
	var status organization_core.OrganizationStatuses
	if r.Status != nil {
		status = organization_core.OrganizationStatuses(*r.Status)
	}

	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	var nameFilter *core.ComparableFilter[string] = nil
	if r.Name != nil {
		nameFilter = &core.ComparableFilter[string]{
			Like: r.Name,
		}
	}

	var statusFilter *core.ComparableFilter[organization_core.OrganizationStatuses] = nil
	if r.Status != nil {
		statusFilter = &core.ComparableFilter[organization_core.OrganizationStatuses]{
			Equals: &status,
		}
	}

	var relationsInput core.RelationsInput = make([]string, 0)
	if r.Relations != nil {
		relationsInput = strings.Split(stringutils.CamelCaseToPascalCase(*r.Relations), ",")
	}

	return organization_services.ListOrganizationsInput{
		Filters: organization_repositories.OrganizationFilters{
			Name:   nameFilter,
			Status: statusFilter,
		},
		ShowDeleted: false,
		Pagination: &core.PaginationInput{
			Page:    r.Page,
			PerPage: r.PerPage,
		},
		SortInput: &core.SortInput{
			By:        r.SortBy,
			Direction: &sortDirection,
		},
		RelationsInput: relationsInput,
	}
}
