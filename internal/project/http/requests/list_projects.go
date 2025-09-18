package project_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
	project_services "github.com/gabrielmrtt/taski/internal/project/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListProjectsRequest struct {
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

func (r *ListProjectsRequest) ToInput() project_services.ListProjectsInput {
	var nameFilter *core.ComparableFilter[string] = nil
	if r.Name != nil {
		nameFilter = &core.ComparableFilter[string]{
			Like: r.Name,
		}
	}

	var statusFilter *core.ComparableFilter[project_core.ProjectStatuses] = nil
	if r.Status != nil {
		projectStatus := project_core.ProjectStatuses(*r.Status)
		statusFilter = &core.ComparableFilter[project_core.ProjectStatuses]{
			Equals: &projectStatus,
		}
	}

	var priorityLevelFilter *core.ComparableFilter[project_core.ProjectPriorityLevels] = nil
	if r.PriorityLevel != nil {
		projectPriorityLevel := project_core.ProjectPriorityLevels(*r.PriorityLevel)
		priorityLevelFilter = &core.ComparableFilter[project_core.ProjectPriorityLevels]{
			Equals: &projectPriorityLevel,
		}
	}

	var sortDirection core.SortDirection
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	return project_services.ListProjectsInput{
		Filters: project_repositories.ProjectFilters{
			Name:          nameFilter,
			Status:        statusFilter,
			PriorityLevel: priorityLevelFilter,
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
