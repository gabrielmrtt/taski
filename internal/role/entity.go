package role

import (
	"slices"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/user"
)

type Permission struct {
	Identity    core.Identity
	Name        string
	Description string
	Slug        PermissionSlugs
}

type NewPermissionInput struct {
	Name        string
	Description string
	Slug        PermissionSlugs
}

func NewPermission(input NewPermissionInput) (*Permission, error) {
	return &Permission{
		Identity:    core.NewIdentityWithoutPublic(),
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
	}, nil
}

type Role struct {
	Identity             core.Identity
	Name                 string
	Slug                 string
	Description          string
	Permissions          []Permission
	OrganizationIdentity *core.Identity
	UserCreatorIdentity  *core.Identity
	UserEditorIdentity   *core.Identity
	IsSystemDefault      bool
	core.Timestamps
	DeletedAt *core.DateTime

	Creator *user.User
	Editor  *user.User
}

type NewRoleInput struct {
	Name                 string
	Description          string
	Permissions          []Permission
	OrganizationIdentity *core.Identity
	UserCreatorIdentity  *core.Identity
	IsSystemDefault      bool
}

func NewRole(input NewRoleInput) (*Role, error) {
	now := core.NewDateTime()

	return &Role{
		Identity:             core.NewIdentity(RoleIdentityPrefix),
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
	now := core.NewDateTime()
	r.Timestamps.UpdatedAt = &now
	return nil
}

func (r *Role) ChangeDescription(description string, userEditorIdentity *core.Identity) error {
	if r.Description == description {
		return nil
	}

	r.Description = description
	r.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	r.Timestamps.UpdatedAt = &now
	return nil
}

func (r *Role) HasPermission(permissionSlug PermissionSlugs) bool {
	for _, p := range r.Permissions {
		if p.Slug == permissionSlug {
			return true
		}
	}

	return false
}

func (r *Role) ClearPermissions(userEditorIdentity *core.Identity) {
	r.Permissions = []Permission{}
	r.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	r.Timestamps.UpdatedAt = &now
}

func (r *Role) AddPermission(permission Permission, userEditorIdentity *core.Identity) {
	r.Permissions = append(r.Permissions, permission)
	r.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	r.Timestamps.UpdatedAt = &now
}

func (r *Role) RemovePermission(permission Permission, userEditorIdentity *core.Identity) {
	r.Permissions = slices.DeleteFunc(r.Permissions, func(p Permission) bool {
		return p.Slug == permission.Slug
	})
	r.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	r.Timestamps.UpdatedAt = &now
}

func (r *Role) Delete() {
	now := core.NewDateTime()
	r.DeletedAt = &now
}

func (r *Role) IsDeleted() bool {
	return r.DeletedAt != nil
}
