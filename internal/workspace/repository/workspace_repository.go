package workspacerepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/workspace"
)

type WorkspaceFilters struct {
	OrganizationIdentity      *core.Identity
	AuthenticatedUserIdentity *core.Identity
	Name                      *core.ComparableFilter[string]
	Description               *core.ComparableFilter[string]
	Color                     *core.ComparableFilter[string]
	Status                    *core.ComparableFilter[workspace.WorkspaceStatuses]
	CreatedAt                 *core.ComparableFilter[int64]
	UpdatedAt                 *core.ComparableFilter[int64]
	DeletedAt                 *core.ComparableFilter[int64]
}

type GetWorkspaceByIdentityParams struct {
	WorkspaceIdentity    core.Identity
	OrganizationIdentity *core.Identity
}

type PaginateWorkspacesParams struct {
	Filters     WorkspaceFilters
	SortInput   *core.SortInput
	Pagination  *core.PaginationInput
	ShowDeleted bool
}

type StoreWorkspaceParams struct {
	Workspace *workspace.Workspace
}

type UpdateWorkspaceParams struct {
	Workspace *workspace.Workspace
}

type DeleteWorkspaceParams struct {
	WorkspaceIdentity core.Identity
}

type WorkspaceRepository interface {
	SetTransaction(tx core.Transaction) error

	GetWorkspaceByIdentity(params GetWorkspaceByIdentityParams) (*workspace.Workspace, error)
	PaginateWorkspacesBy(params PaginateWorkspacesParams) (*core.PaginationOutput[workspace.Workspace], error)

	StoreWorkspace(params StoreWorkspaceParams) (*workspace.Workspace, error)
	UpdateWorkspace(params UpdateWorkspaceParams) error
	DeleteWorkspace(params DeleteWorkspaceParams) error
}
