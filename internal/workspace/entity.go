package workspace

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	"github.com/gabrielmrtt/taski/internal/user"
)

type Workspace struct {
	Identity             core.Identity
	Name                 string
	Description          string
	Color                string
	Status               WorkspaceStatuses
	OrganizationIdentity core.Identity
	UserCreatorIdentity  *core.Identity
	UserEditorIdentity   *core.Identity
	Creator              *user.User
	Editor               *user.User
	Organization         *organization.Organization
	Timestamps           core.Timestamps
	DeletedAt            *core.DateTime
}

type NewWorkspaceInput struct {
	Name                 string
	Description          string
	Color                string
	OrganizationIdentity core.Identity
	UserCreatorIdentity  *core.Identity
}

func NewWorkspace(input NewWorkspaceInput) (*Workspace, error) {
	now := core.NewDateTime()

	if _, err := core.NewName(input.Name); err != nil {
		return nil, err
	}

	if _, err := core.NewDescription(input.Description); err != nil {
		return nil, err
	}

	if _, err := core.NewColor(input.Color); err != nil {
		return nil, err
	}

	return &Workspace{
		Identity:             core.NewIdentity(WorkspaceIdentityPrefix),
		Name:                 input.Name,
		Description:          input.Description,
		Color:                input.Color,
		Status:               WorkspaceStatusActive,
		OrganizationIdentity: input.OrganizationIdentity,
		UserCreatorIdentity:  input.UserCreatorIdentity,
		UserEditorIdentity:   nil,
		Timestamps:           core.Timestamps{CreatedAt: &now, UpdatedAt: nil},
		DeletedAt:            nil,
	}, nil
}

func (w *Workspace) ChangeName(name string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewName(name); err != nil {
		return err
	}

	w.Name = name
	w.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	w.Timestamps.UpdatedAt = &now
	return nil
}

func (w *Workspace) ChangeDescription(description string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewDescription(description); err != nil {
		return err
	}

	w.Description = description
	w.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	w.Timestamps.UpdatedAt = &now
	return nil
}

func (w *Workspace) ChangeColor(color string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewColor(color); err != nil {
		return err
	}

	w.Color = color
	w.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	w.Timestamps.UpdatedAt = &now
	return nil
}

func (w *Workspace) ChangeStatus(status WorkspaceStatuses, userEditorIdentity *core.Identity) error {
	if w.Status == status {
		return nil
	}

	w.Status = status
	w.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	w.Timestamps.UpdatedAt = &now
	return nil
}

func (w *Workspace) Delete() {
	now := core.NewDateTime()
	w.DeletedAt = &now
}

func (w *Workspace) IsActive() bool {
	return w.Status == WorkspaceStatusActive
}

func (w *Workspace) IsInactive() bool {
	return w.Status == WorkspaceStatusInactive
}

func (w *Workspace) IsArchived() bool {
	return w.Status == WorkspaceStatusArchived
}

func (w *Workspace) IsDeleted() bool {
	return w.DeletedAt != nil
}

type WorkspaceUser struct {
	WorkspaceIdentity core.Identity
	User              user.User
	Status            WorkspaceUserStatuses
}

type NewWorkspaceUserInput struct {
	WorkspaceIdentity core.Identity
	User              user.User
	Status            WorkspaceUserStatuses
}

func NewWorkspaceUser(input NewWorkspaceUserInput) (*WorkspaceUser, error) {
	return &WorkspaceUser{
		WorkspaceIdentity: input.WorkspaceIdentity,
		User:              input.User,
		Status:            input.Status,
	}, nil
}

func (w *WorkspaceUser) Activate() {
	w.Status = WorkspaceUserStatusActive
}

func (w *WorkspaceUser) Deactivate() {
	w.Status = WorkspaceUserStatusInactive
}

func (w *WorkspaceUser) Invite() {
	w.Status = WorkspaceUserStatusInvited
}

func (w *WorkspaceUser) IsActive() bool {
	return w.Status == WorkspaceUserStatusActive
}

func (w *WorkspaceUser) IsInactive() bool {
	return w.Status == WorkspaceUserStatusInactive
}

func (w *WorkspaceUser) IsInvited() bool {
	return w.Status == WorkspaceUserStatusInvited
}

func (w *WorkspaceUser) AcceptInvitation() {
	w.Status = WorkspaceUserStatusActive
}

func (w *WorkspaceUser) RefuseInvitation() {
	w.Status = WorkspaceUserStatusRefused
}
