package projectrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
)

type ProjectTaskCategoryFilters struct {
	ProjectIdentity *core.Identity
	Name            *core.ComparableFilter[string]
}

type PaginateProjectTaskCategoryParams struct {
	ShowDeleted bool
	Filters     ProjectTaskCategoryFilters
	SortInput   core.SortInput
	Pagination  core.PaginationInput
}

type GetProjectTaskCategoryByIdentityParams struct {
	ProjectTaskCategoryIdentity *core.Identity
	ProjectIdentity             *core.Identity
}

type StoreProjectTaskCategoryParams struct {
	ProjectTaskCategory *project.ProjectTaskCategory
}

type UpdateProjectTaskCategoryParams struct {
	ProjectTaskCategory *project.ProjectTaskCategory
}

type DeleteProjectTaskCategoryParams struct {
	ProjectTaskCategoryIdentity core.Identity
}

type ProjectTaskCategoryRepository interface {
	SetTransaction(tx core.Transaction) error

	GetProjectTaskCategoryByIdentity(params GetProjectTaskCategoryByIdentityParams) (*project.ProjectTaskCategory, error)
	PaginateProjectTaskCategoryBy(params PaginateProjectTaskCategoryParams) (*core.PaginationOutput[project.ProjectTaskCategory], error)

	StoreProjectTaskCategory(params StoreProjectTaskCategoryParams) (*project.ProjectTaskCategory, error)
	UpdateProjectTaskCategory(params UpdateProjectTaskCategoryParams) error
	DeleteProjectTaskCategory(params DeleteProjectTaskCategoryParams) error
}
