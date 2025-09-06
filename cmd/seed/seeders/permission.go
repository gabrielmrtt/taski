package seeders

import (
	"errors"

	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
)

type PermissionSeeder struct {
	PermissionRepository role_core.PermissionRepository
}

func NewPermissionSeeder(permissionRepository role_core.PermissionRepository) *PermissionSeeder {
	return &PermissionSeeder{
		PermissionRepository: permissionRepository,
	}
}

func checkUniquePermissions(permissions []role_core.Permission) error {
	slugsChecked := make(map[string]struct{})

	for _, permission := range permissions {
		if _, ok := slugsChecked[permission.Slug]; ok {
			return errors.New("there are duplicate slugs. unable to continue.")
		}
		slugsChecked[permission.Slug] = struct{}{}
	}

	return nil
}

func (s *PermissionSeeder) Run() error {
	permissions := []role_core.Permission{
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Update organizations",
			Description: "Allow users to update an organization",
			Slug:        "organizations:update",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Attach users to organizations",
			Description: "Allow users to attach users to an organization",
			Slug:        "organizations:users:create",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Detach users from organizations",
			Description: "Allow users to detach users from an organization",
			Slug:        "organizations:users:delete",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Update users in organizations",
			Description: "Allow users to update users in an organization",
			Slug:        "organizations:users:update",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Create roles",
			Description: "Allow users to create roles",
			Slug:        "roles:create",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Update roles",
			Description: "Allow users to update roles",
			Slug:        "roles:update",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Delete roles",
			Description: "Allow users to delete roles",
			Slug:        "roles:delete",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Create projects",
			Description: "Allow users to create projects",
			Slug:        "organizations:projects:create",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "View projects",
			Description: "Allow users to view projects",
			Slug:        "projects:view",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Update projects",
			Description: "Allow users to update projects",
			Slug:        "projects:update",
		},
		{
			Identity:    core.NewIdentity("perm"),
			Name:        "Delete projects",
			Description: "Allow users to delete projects",
			Slug:        "projects:delete",
		},
	}

	err := checkUniquePermissions(permissions)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		existingPermission, err := s.PermissionRepository.GetPermissionBySlug(role_core.GetPermissionBySlugParams{
			Slug: permission.Slug,
		})
		if err != nil {
			return err
		}

		if existingPermission != nil {
			existingPermission.Name = permission.Name
			existingPermission.Description = permission.Description
			existingPermission.Slug = permission.Slug
			err = s.PermissionRepository.UpdatePermission(existingPermission)
			if err != nil {
				return err
			}
		} else {
			_, err = s.PermissionRepository.StorePermission(&permission)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
