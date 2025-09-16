package workspace_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
)

type WorkspaceUserFilters struct {
	WorkspaceIdentity core.Identity
	UserIdentity      *core.Identity
	Status            *core.ComparableFilter[workspace_core.WorkspaceUserStatuses]
}

type GetWorkspaceUserByIdentityParams struct {
	WorkspaceIdentity core.Identity
	UserIdentity      core.Identity
	RelationsInput    core.RelationsInput
}

type GetWorkspaceUsersByUserIdentityParams struct {
	UserIdentity   core.Identity
	RelationsInput *core.RelationsInput
}

type StoreWorkspaceUserParams struct {
	WorkspaceUser *workspace_core.WorkspaceUser
}

type UpdateWorkspaceUserParams struct {
	WorkspaceUser *workspace_core.WorkspaceUser
}

type DeleteWorkspaceUserParams struct {
	WorkspaceIdentity core.Identity
	UserIdentity      core.Identity
}

type DeleteAllByUserIdentityParams struct {
	UserIdentity core.Identity
}

type WorkspaceUserRepository interface {
	SetTransaction(tx core.Transaction) error

	GetWorkspaceUserByIdentity(params GetWorkspaceUserByIdentityParams) (*workspace_core.WorkspaceUser, error)
	GetWorkspaceUsersByUserIdentity(params GetWorkspaceUsersByUserIdentityParams) ([]workspace_core.WorkspaceUser, error)

	StoreWorkspaceUser(params StoreWorkspaceUserParams) (*workspace_core.WorkspaceUser, error)
	UpdateWorkspaceUser(params UpdateWorkspaceUserParams) error
	DeleteWorkspaceUser(params DeleteWorkspaceUserParams) error
	DeleteAllByUserIdentity(params DeleteAllByUserIdentityParams) error
}
