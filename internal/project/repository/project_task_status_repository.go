package projectrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
)

type ProjectTaskStatusFilters struct {
	ProjectIdentity          *core.Identity
	IsDefault                *bool
	ShouldSetTaskToCompleted *bool
	Name                     *core.ComparableFilter[string]
	Order                    *core.ComparableFilter[int8]
}

type GetProjectTaskStatusByIdentityParams struct {
	ProjectTaskStatusIdentity *core.Identity
	ProjectIdentity           *core.Identity
	IsDefault                 *bool
	ShouldSetTaskToCompleted  *bool
	RelationsInput            core.RelationsInput
}

type PaginateProjectTaskStatusesParams struct {
	ShowDeleted    bool
	Filters        ProjectTaskStatusFilters
	SortInput      core.SortInput
	Pagination     core.PaginationInput
	RelationsInput core.RelationsInput
}

type StoreProjectTaskStatusParams struct {
	ProjectTaskStatus *project.ProjectTaskStatus
}

type UpdateProjectTaskStatusParams struct {
	ProjectTaskStatus *project.ProjectTaskStatus
}

type DeleteProjectTaskStatusParams struct {
	ProjectTaskStatusIdentity core.Identity
}

type ListProjectTaskStatusesByParams struct {
	Filters        ProjectTaskStatusFilters
	SortInput      *core.SortInput
	RelationsInput core.RelationsInput
}

type GetLastTaskStatusOrderParams struct {
	ProjectIdentity *core.Identity
}

type GetTaskStatusByOrderParams struct {
	ProjectIdentity *core.Identity
	Order           int8
	RelationsInput  core.RelationsInput
}

type ProjectTaskStatusRepository interface {
	SetTransaction(tx core.Transaction) error

	GetLastTaskStatusOrder(params GetLastTaskStatusOrderParams) (int8, error)
	GetTaskStatusByOrder(params GetTaskStatusByOrderParams) (*project.ProjectTaskStatus, error)

	GetProjectTaskStatusByIdentity(params GetProjectTaskStatusByIdentityParams) (*project.ProjectTaskStatus, error)
	ListProjectTaskStatusesBy(params ListProjectTaskStatusesByParams) ([]project.ProjectTaskStatus, error)
	PaginateProjectTaskStatusesBy(params PaginateProjectTaskStatusesParams) (*core.PaginationOutput[project.ProjectTaskStatus], error)

	StoreProjectTaskStatus(params StoreProjectTaskStatusParams) (*project.ProjectTaskStatus, error)
	UpdateProjectTaskStatus(params UpdateProjectTaskStatusParams) error
	DeleteProjectTaskStatus(params DeleteProjectTaskStatusParams) error
}
