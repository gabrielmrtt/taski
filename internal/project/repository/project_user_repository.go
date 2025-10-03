package projectrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
)

type ProjectUserFilters struct {
	ProjectIdentity core.Identity
	UserIdentity    *core.Identity
	Status          *core.ComparableFilter[project.ProjectUserStatuses]
}

type GetProjectUserByIdentityParams struct {
	ProjectIdentity core.Identity
	UserIdentity    core.Identity
	RelationsInput  core.RelationsInput
}

type GetProjectUsersByUserIdentityParams struct {
	UserIdentity   core.Identity
	RelationsInput core.RelationsInput
}

type PaginateProjectUsersParams struct {
	Filters        ProjectUserFilters
	SortInput      *core.SortInput
	Pagination     *core.PaginationInput
	RelationsInput core.RelationsInput
}

type StoreProjectUserParams struct {
	ProjectUser *project.ProjectUser
}

type UpdateProjectUserParams struct {
	ProjectUser *project.ProjectUser
}

type DeleteProjectUserParams struct {
	ProjectIdentity core.Identity
	UserIdentity    core.Identity
}

type DeleteAllByUserIdentityParams struct {
	UserIdentity core.Identity
}

type ProjectUserRepository interface {
	SetTransaction(tx core.Transaction) error

	GetProjectUserByIdentity(params GetProjectUserByIdentityParams) (*project.ProjectUser, error)
	GetProjectUsersByUserIdentity(params GetProjectUsersByUserIdentityParams) ([]project.ProjectUser, error)

	StoreProjectUser(params StoreProjectUserParams) (*project.ProjectUser, error)
	UpdateProjectUser(params UpdateProjectUserParams) error
	DeleteProjectUser(params DeleteProjectUserParams) error
	DeleteAllByUserIdentity(params DeleteAllByUserIdentityParams) error
}
