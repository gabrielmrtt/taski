package seeders

import (
	"slices"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/role"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type RolesSeeder struct {
	RoleRepository       rolerepo.RoleRepository
	PermissionRepository rolerepo.PermissionRepository
}

func NewRolesSeeder(roleRepository rolerepo.RoleRepository, permissionRepository rolerepo.PermissionRepository) *RolesSeeder {
	return &RolesSeeder{
		RoleRepository:       roleRepository,
		PermissionRepository: permissionRepository,
	}
}

func (s *RolesSeeder) Run() error {
	permissions, err := s.PermissionRepository.ListPermissionsBy(rolerepo.ListPermissionsParams{
		Filters: rolerepo.PermissionFilters{},
	})
	if err != nil {
		return err
	}

	now := datetimeutils.EpochNow()

	var roles []role.Role = make([]role.Role, 0)
	for _, defaultRole := range role.DefaultRoleSlugsArray {
		var defaultRolePermissions []role.Permission = make([]role.Permission, 0)
		for _, p := range *permissions {
			if slices.Contains(defaultRole.Permissions, p.Slug) {
				defaultRolePermissions = append(defaultRolePermissions, p)
			}
		}

		roles = append(roles, role.Role{
			Identity:        core.NewIdentity(role.RoleIdentityPrefix),
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

	for _, rol := range roles {
		exists, err := s.RoleRepository.GetSystemDefaultRole(rolerepo.GetDefaultRoleParams{
			Slug: role.DefaultRoleSlugs(rol.Slug),
		})
		if err != nil {
			return err
		}
		if exists != nil {
			exists.Name = rol.Name
			exists.Description = rol.Description
			exists.Permissions = rol.Permissions
			exists.UserCreatorIdentity = rol.UserCreatorIdentity
			exists.UserEditorIdentity = rol.UserEditorIdentity
			exists.OrganizationIdentity = rol.OrganizationIdentity
			exists.DeletedAt = rol.DeletedAt
			exists.Timestamps = rol.Timestamps

			err = s.RoleRepository.UpdateRole(rolerepo.UpdateRoleParams{Role: exists})
			if err != nil {
				return err
			}
		} else {
			_, err := s.RoleRepository.StoreRole(rolerepo.StoreRoleParams{Role: &rol})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
