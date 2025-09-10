package organization_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_services "github.com/gabrielmrtt/taski/internal/organization/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListOrganizationUsersRequest struct {
	Name          *string `schema:"name"`
	Email         *string `schema:"email"`
	DisplayName   *string `schema:"display_name"`
	RoleId        *string `schema:"role_id"`
	Status        *string `schema:"status"`
	Page          *int    `schema:"page"`
	PerPage       *int    `schema:"per_page"`
	SortBy        *string `schema:"sort_by"`
	SortDirection *string `schema:"sort_direction"`
}

func (r *ListOrganizationUsersRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListOrganizationUsersRequest) ToInput() organization_services.ListOrganizationUsersInput {
	var status organization_core.OrganizationUserStatuses
	if r.Status != nil {
		status = organization_core.OrganizationUserStatuses(*r.Status)
	}

	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	var nameFilter *core.ComparableFilter[string] = nil
	if r.Name != nil {
		nameFilter = &core.ComparableFilter[string]{
			Equals: r.Name,
		}
	}

	var emailFilter *core.ComparableFilter[string] = nil
	if r.Email != nil {
		emailFilter = &core.ComparableFilter[string]{
			Equals: r.Email,
		}
	}

	var displayNameFilter *core.ComparableFilter[string] = nil
	if r.DisplayName != nil {
		displayNameFilter = &core.ComparableFilter[string]{
			Equals: r.DisplayName,
		}
	}

	var roleIdFilter *core.ComparableFilter[string] = nil
	if r.RoleId != nil {
		roleIdFilter = &core.ComparableFilter[string]{
			Equals: r.RoleId,
		}
	}

	var statusFilter *core.ComparableFilter[organization_core.OrganizationUserStatuses] = nil
	if r.Status != nil {
		statusFilter = &core.ComparableFilter[organization_core.OrganizationUserStatuses]{
			Equals: &status,
		}
	}

	return organization_services.ListOrganizationUsersInput{
		Filters: organization_core.OrganizationUserFilters{
			Name:         nameFilter,
			Email:        emailFilter,
			DisplayName:  displayNameFilter,
			RolePublicId: roleIdFilter,
			Status:       statusFilter,
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
