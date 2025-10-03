package seeders

import (
	"errors"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/role"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
)

type PermissionSeeder struct {
	PermissionRepository rolerepo.PermissionRepository
}

func NewPermissionSeeder(permissionRepository rolerepo.PermissionRepository) *PermissionSeeder {
	return &PermissionSeeder{
		PermissionRepository: permissionRepository,
	}
}

func checkUniquePermissions(permissions []role.Permission) error {
	slugsChecked := make(map[string]struct{})

	for _, permission := range permissions {
		if _, ok := slugsChecked[string(permission.Slug)]; ok {
			return errors.New("there are duplicate slugs. unable to continue")
		}
		slugsChecked[string(permission.Slug)] = struct{}{}
	}

	return nil
}

func (s *PermissionSeeder) Run() error {
	var permissions []role.Permission = make([]role.Permission, 0)

	for _, i := range role.PermissionSlugsArray {
		permissions = append(permissions, role.Permission{
			Identity:    core.NewIdentityWithoutPublic(),
			Name:        i.Name,
			Description: i.Description,
			Slug:        i.Slug,
		})
	}

	err := checkUniquePermissions(permissions)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		existingPermission, err := s.PermissionRepository.GetPermissionBySlug(rolerepo.GetPermissionBySlugParams{
			Slug: string(permission.Slug),
		})

		if err != nil {
			return err
		}

		if existingPermission != nil {
			existingPermission.Name = permission.Name
			existingPermission.Description = permission.Description
			err = s.PermissionRepository.UpdatePermission(rolerepo.UpdatePermissionParams{Permission: existingPermission})
			if err != nil {
				return err
			}
		} else {
			_, err = s.PermissionRepository.StorePermission(rolerepo.StorePermissionParams{Permission: &permission})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
