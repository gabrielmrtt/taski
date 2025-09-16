package workspace_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
)

type WorkspaceFilters struct {
	OrganizationIdentity core.Identity
	Name                 *core.ComparableFilter[string]
	Description          *core.ComparableFilter[string]
	Color                *core.ComparableFilter[string]
	Status               *core.ComparableFilter[workspace_core.WorkspaceStatuses]
	CreatedAt            *core.ComparableFilter[int64]
	UpdatedAt            *core.ComparableFilter[int64]
	DeletedAt            *core.ComparableFilter[int64]
}

type GetWorkspaceByIdentityParams struct {
	WorkspaceIdentity    core.Identity
	OrganizationIdentity *core.Identity
	RelationsInput       *core.RelationsInput
}

type PaginateWorkspacesParams struct {
	Filters        WorkspaceFilters
	SortInput      *core.SortInput
	Pagination     *core.PaginationInput
	RelationsInput *core.RelationsInput
}

type StoreWorkspaceParams struct {
	Workspace *workspace_core.Workspace
}

type UpdateWorkspaceParams struct {
	Workspace *workspace_core.Workspace
}

type DeleteWorkspaceParams struct {
	WorkspaceIdentity core.Identity
}

type WorkspaceRepository interface {
	SetTransaction(tx core.Transaction) error

	GetWorkspaceByIdentity(params GetWorkspaceByIdentityParams) (*workspace_core.Workspace, error)
	PaginateWorkspacesBy(params PaginateWorkspacesParams) (*core.PaginationOutput[workspace_core.Workspace], error)

	StoreWorkspace(params StoreWorkspaceParams) (*workspace_core.Workspace, error)
	UpdateWorkspace(params UpdateWorkspaceParams) error
	DeleteWorkspace(params DeleteWorkspaceParams) error
}
