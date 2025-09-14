package seeders

import (
	"errors"
	"fmt"

	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
)

type PermissionSeeder struct {
	PermissionRepository role_repositories.PermissionRepository
}

func NewPermissionSeeder(permissionRepository role_repositories.PermissionRepository) *PermissionSeeder {
	return &PermissionSeeder{
		PermissionRepository: permissionRepository,
	}
}

func checkUniquePermissions(permissions []role_core.Permission) error {
	slugsChecked := make(map[string]struct{})

	for _, permission := range permissions {
		if _, ok := slugsChecked[string(permission.Slug)]; ok {
			return errors.New("there are duplicate slugs. unable to continue.")
		}
		slugsChecked[string(permission.Slug)] = struct{}{}
	}

	return nil
}

func (s *PermissionSeeder) Run() error {
	var permissions []role_core.Permission = make([]role_core.Permission, 0)

	for _, i := range role_core.PermissionSlugsArray {
		permissions = append(permissions, role_core.Permission{
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
		existingPermission, err := s.PermissionRepository.GetPermissionBySlug(role_repositories.GetPermissionBySlugParams{
			Slug: string(permission.Slug),
		})

		if err != nil {
			return err
		}

		if existingPermission != nil {
			existingPermission.Name = permission.Name
			existingPermission.Description = permission.Description
			err = s.PermissionRepository.UpdatePermission(role_repositories.UpdatePermissionParams{Permission: existingPermission})
			if err != nil {
				fmt.Println("caiu aqui")
				return err
			}
		} else {
			_, err = s.PermissionRepository.StorePermission(role_repositories.StorePermissionParams{Permission: &permission})
			if err != nil {
				fmt.Println("caiu aqui 2")
				return err
			}
		}
	}

	return nil
}
