package workspace_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
	workspace_services "github.com/gabrielmrtt/taski/internal/workspace/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListWorkspacesRequest struct {
	Name          *string `json:"name" schema:"name"`
	Description   *string `json:"description" schema:"description"`
	Status        *string `json:"status" schema:"status"`
	Page          *int    `json:"page" schema:"page"`
	PerPage       *int    `json:"perPage" schema:"perPage"`
	SortBy        *string `json:"sortBy" schema:"sortBy"`
	SortDirection *string `json:"sortDirection" schema:"sortDirection"`
}

func (r *ListWorkspacesRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListWorkspacesRequest) ToInput() workspace_services.ListWorkspacesInput {
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

	var descriptionFilter *core.ComparableFilter[string] = nil
	if r.Description != nil {
		descriptionFilter = &core.ComparableFilter[string]{
			Like: r.Description,
		}
	}

	var statusFilter *core.ComparableFilter[workspace_core.WorkspaceStatuses] = nil
	if r.Status != nil {
		workspaceStatus := workspace_core.WorkspaceStatuses(*r.Status)
		statusFilter = &core.ComparableFilter[workspace_core.WorkspaceStatuses]{
			Equals: &workspaceStatus,
		}
	}

	return workspace_services.ListWorkspacesInput{
		Filters: workspace_repositories.WorkspaceFilters{
			Name:        nameFilter,
			Description: descriptionFilter,
			Status:      statusFilter,
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
