package organizationhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	organizationservice "github.com/gabrielmrtt/taski/internal/organization/service"
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

func (r *ListOrganizationsRequest) ToInput() organizationservice.ListOrganizationsInput {
	var status organization.OrganizationStatuses
	if r.Status != nil {
		status = organization.OrganizationStatuses(*r.Status)
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

	var statusFilter *core.ComparableFilter[organization.OrganizationStatuses] = nil
	if r.Status != nil {
		statusFilter = &core.ComparableFilter[organization.OrganizationStatuses]{
			Equals: &status,
		}
	}

	return organizationservice.ListOrganizationsInput{
		Filters: organizationrepo.OrganizationFilters{
			Name:   nameFilter,
			Status: statusFilter,
		},
		Pagination: core.PaginationInput{
			Page:    r.Page,
			PerPage: r.PerPage,
		},
		SortInput: core.SortInput{
			By:        r.SortBy,
			Direction: &sortDirection,
		},
		RelationsInput: corehttp.GetRelationsInput(*r.Relations),
	}
}
