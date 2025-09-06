package role_core

import (
	"slices"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type Permission struct {
	Identity    core.Identity
	Name        string
	Description string
	Slug        string
}

type NewPermissionInput struct {
	Name        string
	Description string
	Slug        string
}

func NewPermission(input NewPermissionInput) (*Permission, error) {
	return &Permission{
		Identity:    core.NewIdentity("perm"),
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
	}, nil
}

type Role struct {
	Identity             core.Identity
	Name                 string
	Description          string
	Permissions          []Permission
	OrganizationIdentity core.Identity
	UserCreatorIdentity  *core.Identity
	UserEditorIdentity   *core.Identity
	IsSystemDefault      bool
	core.Timestamps
	DeletedAt *int64
}

type NewRoleInput struct {
	Name                 string
	Description          string
	Permissions          []Permission
	OrganizationIdentity core.Identity
	UserCreatorIdentity  *core.Identity
	IsSystemDefault      bool
}

func NewRole(input NewRoleInput) (*Role, error) {
	now := datetimeutils.EpochNow()

	return &Role{
		Identity:             core.NewIdentity("role"),
		Name:                 input.Name,
		Description:          input.Description,
		Permissions:          input.Permissions,
		OrganizationIdentity: input.OrganizationIdentity,
		UserCreatorIdentity:  input.UserCreatorIdentity,
		UserEditorIdentity:   nil,
		IsSystemDefault:      input.IsSystemDefault,
		Timestamps: core.Timestamps{
			CreatedAt: &now,
			UpdatedAt: nil,
		},
	}, nil
}

func (r *Role) ChangeName(name string, userEditorIdentity *core.Identity) error {
	if r.Name == name {
		return nil
	}

	r.Name = name
	r.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	r.Timestamps.UpdatedAt = &now
	return nil
}

func (r *Role) ChangeDescription(description string, userEditorIdentity *core.Identity) error {
	if r.Description == description {
		return nil
	}

	r.Description = description
	r.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	r.Timestamps.UpdatedAt = &now
	return nil
}

func (r *Role) HasPermission(permission Permission) bool {
	return slices.Contains(r.Permissions, permission)
}

func (r *Role) ClearPermissions(userEditorIdentity *core.Identity) {
	r.Permissions = []Permission{}
	r.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	r.Timestamps.UpdatedAt = &now
}

func (r *Role) AddPermission(permission Permission, userEditorIdentity *core.Identity) {
	r.Permissions = append(r.Permissions, permission)
	r.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	r.Timestamps.UpdatedAt = &now
}

func (r *Role) RemovePermission(permission Permission, userEditorIdentity *core.Identity) {
	r.Permissions = slices.DeleteFunc(r.Permissions, func(p Permission) bool {
		return p.Slug == permission.Slug
	})
	r.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	r.Timestamps.UpdatedAt = &now
}

func (r *Role) Delete() {
	now := datetimeutils.EpochNow()
	r.DeletedAt = &now
}

func (r *Role) IsDeleted() bool {
	return r.DeletedAt != nil
}
