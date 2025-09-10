package seeders

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type RolesSeeder struct {
	RoleRepository       role_core.RoleRepository
	PermissionRepository role_core.PermissionRepository
}

func NewRolesSeeder(roleRepository role_core.RoleRepository, permissionRepository role_core.PermissionRepository) *RolesSeeder {
	return &RolesSeeder{
		RoleRepository:       roleRepository,
		PermissionRepository: permissionRepository,
	}
}

func (s *RolesSeeder) Run() error {
	permissions, err := s.PermissionRepository.ListPermissionsBy(role_core.ListPermissionsParams{
		Filters: role_core.PermissionFilters{},
		Include: map[string]any{},
	})
	if err != nil {
		return err
	}

	now := datetimeutils.EpochNow()

	roles := []role_core.Role{
		{
			Identity:        core.NewIdentity("role"),
			Name:            "Admin",
			Slug:            "admin",
			IsSystemDefault: true,
			Description:     "Admin role",
			Permissions:     *permissions,
			Timestamps: core.Timestamps{
				CreatedAt: &now,
				UpdatedAt: nil,
			},
			UserCreatorIdentity:  nil,
			UserEditorIdentity:   nil,
			OrganizationIdentity: nil,
			DeletedAt:            nil,
		},
		{
			Identity:        core.NewIdentity("role"),
			Name:            "Default",
			Slug:            "default",
			IsSystemDefault: true,
			Description:     "Default role",
			Permissions:     *permissions,
			Timestamps: core.Timestamps{
				CreatedAt: &now,
				UpdatedAt: nil,
			},
			UserCreatorIdentity:  nil,
			UserEditorIdentity:   nil,
			OrganizationIdentity: nil,
			DeletedAt:            nil,
		},
	}

	for _, role := range roles {
		exists, err := s.RoleRepository.GetSystemDefaultRole(role_core.GetDefaultRoleParams{
			Slug: role.Slug,
		})
		if err != nil {
			return err
		}
		if exists != nil {
			exists.Name = role.Name
			exists.Description = role.Description
			exists.Permissions = role.Permissions
			exists.UserCreatorIdentity = role.UserCreatorIdentity
			exists.UserEditorIdentity = role.UserEditorIdentity
			exists.OrganizationIdentity = role.OrganizationIdentity
			exists.DeletedAt = role.DeletedAt
			exists.Timestamps = role.Timestamps

			err = s.RoleRepository.UpdateRole(exists)
			if err != nil {
				return err
			}
		} else {
			_, err := s.RoleRepository.StoreRole(&role)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
