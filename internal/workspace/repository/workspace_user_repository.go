package workspacerepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/workspace"
)

type WorkspaceUserFilters struct {
	WorkspaceIdentity core.Identity
	UserIdentity      *core.Identity
	Status            *core.ComparableFilter[workspace.WorkspaceUserStatuses]
}

type GetWorkspaceUserByIdentityParams struct {
	WorkspaceIdentity core.Identity
	UserIdentity      core.Identity
	RelationsInput    core.RelationsInput
}

type GetWorkspaceUsersByUserIdentityParams struct {
	UserIdentity   core.Identity
	RelationsInput core.RelationsInput
}

type StoreWorkspaceUserParams struct {
	WorkspaceUser *workspace.WorkspaceUser
}

type UpdateWorkspaceUserParams struct {
	WorkspaceUser *workspace.WorkspaceUser
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

	GetWorkspaceUserByIdentity(params GetWorkspaceUserByIdentityParams) (*workspace.WorkspaceUser, error)
	GetWorkspaceUsersByUserIdentity(params GetWorkspaceUsersByUserIdentityParams) ([]workspace.WorkspaceUser, error)

	StoreWorkspaceUser(params StoreWorkspaceUserParams) (*workspace.WorkspaceUser, error)
	UpdateWorkspaceUser(params UpdateWorkspaceUserParams) error
	DeleteWorkspaceUser(params DeleteWorkspaceUserParams) error
	DeleteAllByUserIdentity(params DeleteAllByUserIdentityParams) error
}
