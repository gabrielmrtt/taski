package project_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project_core "github.com/gabrielmrtt/taski/internal/project"
)

type ProjectFilters struct {
	WorkspaceIdentity core.Identity
	Name              *core.ComparableFilter[string]
	Description       *core.ComparableFilter[string]
	Color             *core.ComparableFilter[string]
	PriorityLevel     *core.ComparableFilter[project_core.ProjectPriorityLevels]
	Status            *core.ComparableFilter[project_core.ProjectStatuses]
	CreatedAt         *core.ComparableFilter[int64]
	UpdatedAt         *core.ComparableFilter[int64]
	DeletedAt         *core.ComparableFilter[int64]
}

type GetProjectByIdentityParams struct {
	ProjectIdentity   core.Identity
	WorkspaceIdentity *core.Identity
	RelationsInput    *core.RelationsInput
}

type PaginateProjectsParams struct {
	Filters        ProjectFilters
	SortInput      *core.SortInput
	Pagination     *core.PaginationInput
	RelationsInput *core.RelationsInput
}

type StoreProjectParams struct {
	Project *project_core.Project
}

type UpdateProjectParams struct {
	Project *project_core.Project
}

type DeleteProjectParams struct {
	ProjectIdentity core.Identity
}

type ProjectRepository interface {
	SetTransaction(tx core.Transaction) error

	GetProjectByIdentity(params GetProjectByIdentityParams) (*project_core.Project, error)
	PaginateProjectsBy(params PaginateProjectsParams) (*core.PaginationOutput[project_core.Project], error)

	StoreProject(params StoreProjectParams) (*project_core.Project, error)
	UpdateProject(params UpdateProjectParams) error
	DeleteProject(params DeleteProjectParams) error
}
