package projectrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
)

type ProjectFilters struct {
	OrganizationIdentity      *core.Identity
	WorkspaceIdentity         *core.Identity
	AuthenticatedUserIdentity *core.Identity
	Name                      *core.ComparableFilter[string]
	Description               *core.ComparableFilter[string]
	Color                     *core.ComparableFilter[string]
	PriorityLevel             *core.ComparableFilter[project.ProjectPriorityLevels]
	Status                    *core.ComparableFilter[project.ProjectStatuses]
	CreatedAt                 *core.ComparableFilter[int64]
	UpdatedAt                 *core.ComparableFilter[int64]
	DeletedAt                 *core.ComparableFilter[int64]
}

type GetProjectByIdentityParams struct {
	ProjectIdentity      core.Identity
	WorkspaceIdentity    *core.Identity
	OrganizationIdentity *core.Identity
	RelationsInput       core.RelationsInput
}

type PaginateProjectsParams struct {
	Filters        ProjectFilters
	SortInput      *core.SortInput
	Pagination     *core.PaginationInput
	RelationsInput core.RelationsInput
}

type StoreProjectParams struct {
	Project *project.Project
}

type UpdateProjectParams struct {
	Project *project.Project
}

type DeleteProjectParams struct {
	ProjectIdentity core.Identity
}

type ProjectRepository interface {
	SetTransaction(tx core.Transaction) error

	GetProjectByIdentity(params GetProjectByIdentityParams) (*project.Project, error)
	PaginateProjectsBy(params PaginateProjectsParams) (*core.PaginationOutput[project.Project], error)

	StoreProject(params StoreProjectParams) (*project.Project, error)
	UpdateProject(params UpdateProjectParams) error
	DeleteProject(params DeleteProjectParams) error
}
