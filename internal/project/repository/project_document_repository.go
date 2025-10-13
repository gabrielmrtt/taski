package projectrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
)

type ProjectDocumentVersionManagerFilters struct {
	ProjectIdentity *core.Identity
	Title           *core.ComparableFilter[string]
	CreatedAt       *core.ComparableFilter[int64]
}

type PaginateProjectDocumentVersionManagersByParams struct {
	Filters    ProjectDocumentVersionManagerFilters
	SortInput  core.SortInput
	Pagination core.PaginationInput
}

type ProjectDocumentVersionFilters struct {
	ProjectDocumentVersionManagerIdentity *core.Identity
	Version                               *core.ComparableFilter[string]
	CreatedAt                             *core.ComparableFilter[int64]
}

type PaginateProjectDocumentVersionsByParams struct {
	Filters    ProjectDocumentVersionFilters
	SortInput  core.SortInput
	Pagination core.PaginationInput
}

type GetProjectDocumentVersionByParams struct {
	ProjectDocumentVersionManagerIdentity *core.Identity
	ProjectDocumentVersionIdentity        core.Identity
}

type GetProjectDocumentVersionManagerByParams struct {
	ProjectDocumentVersionManagerIdentity core.Identity
}

type StoreProjectDocumentVersionManagerParams struct {
	ProjectDocumentVersionManager *project.ProjectDocumentVersionManager
}

type StoreProjectDocumentVersionParams struct {
	ProjectDocumentVersion *project.ProjectDocumentVersion
}

type UpdateProjectDocumentVersionParams struct {
	ProjectDocumentVersion *project.ProjectDocumentVersion
}

type DeleteProjectDocumentVersionParams struct {
	ProjectDocumentVersionIdentity core.Identity
}

type DeleteProjectDocumentVersionManagerParams struct {
	ProjectDocumentVersionManagerIdentity core.Identity
}

type ListProjectDocumentVersionsByProjectDocumentVersionManagerIdentityParams struct {
	ProjectDocumentVersionManagerIdentity core.Identity
	SortInput                             core.SortInput
}

type ProjectDocumentRepository interface {
	SetTransaction(tx core.Transaction) error

	PaginateProjectDocumentVersionManagersBy(params PaginateProjectDocumentVersionManagersByParams) (*core.PaginationOutput[project.ProjectDocumentVersionManager], error)
	PaginateProjectDocumentVersionsBy(params PaginateProjectDocumentVersionsByParams) (*core.PaginationOutput[project.ProjectDocumentVersion], error)
	ListProjectDocumentVersionsByProjectDocumentVersionManagerIdentity(params ListProjectDocumentVersionsByProjectDocumentVersionManagerIdentityParams) ([]project.ProjectDocumentVersion, error)
	GetProjectDocumentVersionBy(params GetProjectDocumentVersionByParams) (*project.ProjectDocumentVersion, error)

	StoreProjectDocumentVersionManager(params StoreProjectDocumentVersionManagerParams) (*project.ProjectDocumentVersionManager, error)
	GetProjectDocumentVersionManagerBy(params GetProjectDocumentVersionManagerByParams) (*project.ProjectDocumentVersionManager, error)
	DeleteProjectDocumentVersionManager(params DeleteProjectDocumentVersionManagerParams) error

	StoreProjectDocumentVersion(params StoreProjectDocumentVersionParams) (*project.ProjectDocumentVersion, error)
	UpdateProjectDocumentVersion(params UpdateProjectDocumentVersionParams) error
	DeleteProjectDocumentVersion(params DeleteProjectDocumentVersionParams) error
}
