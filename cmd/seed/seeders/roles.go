package seeders

import (
	"slices"

	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type RolesSeeder struct {
	RoleRepository       role_repositories.RoleRepository
	PermissionRepository role_repositories.PermissionRepository
}

func NewRolesSeeder(roleRepository role_repositories.RoleRepository, permissionRepository role_repositories.PermissionRepository) *RolesSeeder {
	return &RolesSeeder{
		RoleRepository:       roleRepository,
		PermissionRepository: permissionRepository,
	}
}

func (s *RolesSeeder) Run() error {
	permissions, err := s.PermissionRepository.ListPermissionsBy(role_repositories.ListPermissionsParams{
		Filters: role_repositories.PermissionFilters{},
	})
	if err != nil {
		return err
	}

	now := datetimeutils.EpochNow()

	var roles []role_core.Role = make([]role_core.Role, 0)
	for _, defaultRole := range role_core.DefaultRoleSlugsArray {
		var defaultRolePermissions []role_core.Permission = make([]role_core.Permission, 0)
		for _, p := range *permissions {
			if slices.Contains(defaultRole.Permissions, p.Slug) {
				defaultRolePermissions = append(defaultRolePermissions, p)
			}
		}

		roles = append(roles, role_core.Role{
			Identity:        core.NewIdentity(role_core.RoleIdentityPrefix),
			Name:            defaultRole.Name,
			Slug:            string(defaultRole.Slug),
			IsSystemDefault: true,
			Description:     defaultRole.Description,
			Permissions:     defaultRolePermissions,
			Timestamps: core.Timestamps{
				CreatedAt: &now,
				UpdatedAt: nil,
			},
			UserCreatorIdentity:  nil,
			UserEditorIdentity:   nil,
			OrganizationIdentity: nil,
			DeletedAt:            nil,
		})
	}

	for _, role := range roles {
		exists, err := s.RoleRepository.GetSystemDefaultRole(role_repositories.GetDefaultRoleParams{
			Slug: role_core.DefaultRoleSlugs(role.Slug),
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

			err = s.RoleRepository.UpdateRole(role_repositories.UpdateRoleParams{Role: exists})
			if err != nil {
				return err
			}
		} else {
			_, err := s.RoleRepository.StoreRole(role_repositories.StoreRoleParams{Role: &role})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
