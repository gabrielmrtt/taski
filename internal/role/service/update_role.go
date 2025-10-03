package roleservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
)

type UpdateRoleService struct {
	RoleRepository        rolerepo.RoleRepository
	PermissionRepository  rolerepo.PermissionRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateRoleService(
	roleRepository rolerepo.RoleRepository,
	permissionRepository rolerepo.PermissionRepository,
	transactionRepository core.TransactionRepository,
) *UpdateRoleService {
	return &UpdateRoleService{
		RoleRepository:        roleRepository,
		PermissionRepository:  permissionRepository,
		TransactionRepository: transactionRepository,
	}
}

type UpdateRoleInput struct {
	RoleIdentity         core.Identity
	Name                 *string
	Description          *string
	Permissions          []string
	OrganizationIdentity core.Identity
	UserEditorIdentity   core.Identity
}

func (i UpdateRoleInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Name != nil {
		_, err := core.NewName(*i.Name)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "name",
				Error: err.Error(),
			})
		}
	}

	if i.Description != nil {
		_, err := core.NewDescription(*i.Description)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "description",
				Error: err.Error(),
			})
		}
	}

	duplicates := make(map[string]struct{})

	for _, permission := range i.Permissions {
		if _, ok := duplicates[permission]; ok {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "permissions",
				Error: "duplicate permission",
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateRoleService) Execute(input UpdateRoleInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.RoleRepository.SetTransaction(tx)
	s.PermissionRepository.SetTransaction(tx)

	rol, err := s.RoleRepository.GetRoleByIdentityAndOrganizationIdentity(rolerepo.GetRoleByIdentityAndOrganizationIdentityParams{
		RoleIdentity:         input.RoleIdentity,
		OrganizationIdentity: input.OrganizationIdentity,
	})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if rol == nil {
		tx.Rollback()
		return core.NewNotFoundError("role not found")
	}

	if input.Name != nil {
		err = rol.ChangeName(*input.Name, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Description != nil {
		err = rol.ChangeDescription(*input.Description, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Permissions != nil {
		rol.ClearPermissions(&input.UserEditorIdentity)

		for _, permission := range input.Permissions {
			permission, err := s.PermissionRepository.GetPermissionBySlug(rolerepo.GetPermissionBySlugParams{
				Slug: permission,
			})
			if err != nil {
				tx.Rollback()
				return core.NewInternalError(err.Error())
			}

			if permission == nil {
				tx.Rollback()
				return core.NewNotFoundError("permission not found")
			}

			rol.AddPermission(*permission, &input.UserEditorIdentity)
		}
	}

	err = s.RoleRepository.UpdateRole(rolerepo.UpdateRoleParams{Role: rol})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	return nil
}
