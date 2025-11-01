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

type ListOrganizationUsersRequest struct {
	Name          *string `json:"name" schema:"name"`
	Email         *string `json:"email" schema:"email"`
	DisplayName   *string `json:"displayName" schema:"displayName"`
	RoleId        *string `json:"roleId" schema:"roleId"`
	Status        *string `json:"status" schema:"status"`
	Page          *int    `json:"page" schema:"page"`
	PerPage       *int    `json:"perPage" schema:"perPage"`
	SortBy        *string `json:"sortBy" schema:"sortBy"`
	SortDirection *string `json:"sortDirection" schema:"sortDirection"`
	Relations     *string `json:"relations" schema:"relations"`
}

func (r *ListOrganizationUsersRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListOrganizationUsersRequest) ToInput() organizationservice.ListOrganizationUsersInput {
	var status organization.OrganizationUserStatuses
	if r.Status != nil {
		status = organization.OrganizationUserStatuses(*r.Status)
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

	var emailFilter *core.ComparableFilter[string] = nil
	if r.Email != nil {
		emailFilter = &core.ComparableFilter[string]{
			Like: r.Email,
		}
	}

	var displayNameFilter *core.ComparableFilter[string] = nil
	if r.DisplayName != nil {
		displayNameFilter = &core.ComparableFilter[string]{
			Like: r.DisplayName,
		}
	}

	var roleIdFilter *core.ComparableFilter[string] = nil
	if r.RoleId != nil {
		roleIdFilter = &core.ComparableFilter[string]{
			Equals: r.RoleId,
		}
	}

	var statusFilter *core.ComparableFilter[organization.OrganizationUserStatuses] = nil
	if r.Status != nil {
		statusFilter = &core.ComparableFilter[organization.OrganizationUserStatuses]{
			Equals: &status,
		}
	}

	return organizationservice.ListOrganizationUsersInput{
		Filters: organizationrepo.OrganizationUserFilters{
			Name:         nameFilter,
			Email:        emailFilter,
			DisplayName:  displayNameFilter,
			RolePublicId: roleIdFilter,
			Status:       statusFilter,
		},
		Pagination: core.PaginationInput{
			Page:    r.Page,
			PerPage: r.PerPage,
		},
		SortInput: core.SortInput{
			By:        r.SortBy,
			Direction: &sortDirection,
		},
		RelationsInput: corehttp.GetRelationsInput(r.Relations),
	}
}
