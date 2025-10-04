package projecthttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListProjectsRequest struct {
	WorkspaceId   *string `json:"workspaceId"`
	Name          *string `json:"name"`
	Status        *string `json:"status"`
	PriorityLevel *int8   `json:"priorityLevel"`
	Page          *int    `json:"page"`
	PerPage       *int    `json:"perPage"`
	SortBy        *string `json:"sortBy"`
	SortDirection *string `json:"sortDirection"`
}

func (r *ListProjectsRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListProjectsRequest) ToInput() projectservice.ListProjectsInput {
	var workspaceIdentity *core.Identity = nil
	if r.WorkspaceId != nil {
		identity := core.NewIdentityFromPublic(*r.WorkspaceId)
		workspaceIdentity = &identity
	}

	var nameFilter *core.ComparableFilter[string] = nil
	if r.Name != nil {
		nameFilter = &core.ComparableFilter[string]{
			Like: r.Name,
		}
	}

	var statusFilter *core.ComparableFilter[project.ProjectStatuses] = nil
	if r.Status != nil {
		projectStatus := project.ProjectStatuses(*r.Status)
		statusFilter = &core.ComparableFilter[project.ProjectStatuses]{
			Equals: &projectStatus,
		}
	}

	var priorityLevelFilter *core.ComparableFilter[project.ProjectPriorityLevels] = nil
	if r.PriorityLevel != nil {
		projectPriorityLevel := project.ProjectPriorityLevels(*r.PriorityLevel)
		priorityLevelFilter = &core.ComparableFilter[project.ProjectPriorityLevels]{
			Equals: &projectPriorityLevel,
		}
	}

	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	return projectservice.ListProjectsInput{
		Filters: projectrepo.ProjectFilters{
			WorkspaceIdentity: workspaceIdentity,
			Name:              nameFilter,
			Status:            statusFilter,
			PriorityLevel:     priorityLevelFilter,
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
