package organization_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListOrganizationsRequest struct {
	Name          *string `schema:"name"`
	Status        *string `schema:"status"`
	Page          *int    `schema:"page"`
	PerPage       *int    `schema:"per_page"`
	SortBy        *string `schema:"sort_by"`
	SortDirection *string `schema:"sort_direction"`
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

	return organization_services.ListOrganizationsInput{
		Filters: organization_core.OrganizationFilters{
			Name: &core.ComparableFilter[string]{
				Equals: r.Name,
			},
			Status: &core.ComparableFilter[organization_core.OrganizationStatuses]{
				Equals: &status,
			},
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
	}
}
