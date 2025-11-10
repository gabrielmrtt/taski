package taskhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type ListTasksRequest struct {
	ProjectId      *string `json:"projectId"`
	StatusId       *string `json:"statusId"`
	CategoryId     *string `json:"categoryId"`
	ParentTaskId   *string `json:"parentTaskId"`
	Name           *string `json:"name"`
	Completed      *bool   `json:"completed"`
	CompletedAtLte *int64  `json:"completedAtLte"`
	CompletedAtGte *int64  `json:"completedAtGte"`
	DueDateLte     *int64  `json:"dueDateLte"`
	DueDateGte     *int64  `json:"dueDateGte"`
	Page           *int    `json:"page"`
	PerPage        *int    `json:"perPage"`
	SortBy         *string `json:"sortBy"`
	SortDirection  *string `json:"sortDirection"`
	Relations      *string `json:"relations"`
}

func (r *ListTasksRequest) FromQuery(ctx *gin.Context) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(r, ctx.Request.URL.Query())
}

func (r *ListTasksRequest) ToInput() taskservice.ListTasksInput {
	var projectIdentity *core.Identity = nil
	if r.ProjectId != nil {
		identity := core.NewIdentity(*r.ProjectId)
		projectIdentity = &identity
	}

	var statusIdentity *core.Identity = nil
	if r.StatusId != nil {
		identity := core.NewIdentity(*r.StatusId)
		statusIdentity = &identity
	}

	var categoryIdentity *core.Identity = nil
	if r.CategoryId != nil {
		identity := core.NewIdentity(*r.CategoryId)
		categoryIdentity = &identity
	}

	var parentTaskIdentity *core.Identity = nil
	if r.ParentTaskId != nil {
		identity := core.NewIdentity(*r.ParentTaskId)
		parentTaskIdentity = &identity
	}

	var nameFilter *core.ComparableFilter[string] = nil
	if r.Name != nil {
		nameFilter = &core.ComparableFilter[string]{
			Like: r.Name,
		}
	}

	var completedAtFilter *core.ComparableFilter[int64] = &core.ComparableFilter[int64]{}
	if r.Completed != nil {
		completedAtFilter.NotNull = r.Completed
	}

	if r.CompletedAtLte != nil {
		completedAtFilter.LessThanOrEqual = r.CompletedAtLte
	}

	if r.CompletedAtGte != nil {
		completedAtFilter.GreaterThanOrEqual = r.CompletedAtGte
	}

	var dueDateFilter *core.ComparableFilter[int64] = &core.ComparableFilter[int64]{}
	if r.DueDateLte != nil {
		dueDateFilter.LessThanOrEqual = r.DueDateLte
	}

	if r.DueDateGte != nil {
		dueDateFilter.GreaterThanOrEqual = r.DueDateGte
	}

	var sortDirection core.SortDirection = core.SortDirectionAsc
	if r.SortDirection != nil {
		sortDirection = core.SortDirection(*r.SortDirection)
	}

	return taskservice.ListTasksInput{
		Filters: taskrepo.TaskFilters{
			ProjectIdentity:      projectIdentity,
			TaskStatusIdentity:   statusIdentity,
			TaskCategoryIdentity: categoryIdentity,
			ParentTaskIdentity:   parentTaskIdentity,
			Name:                 nameFilter,
			CompletedAt:          completedAtFilter,
			DueDate:              dueDateFilter,
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
