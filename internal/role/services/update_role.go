package role_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
)

type UpdateRoleService struct {
	RoleRepository        role_repositories.RoleRepository
	PermissionRepository  role_repositories.PermissionRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateRoleService(
	roleRepository role_repositories.RoleRepository,
	permissionRepository role_repositories.PermissionRepository,
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

	organizationHasUser, err := s.RoleRepository.CheckIfOrganizatonHasUser(input.OrganizationIdentity, input.UserEditorIdentity)
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if !organizationHasUser {
		tx.Rollback()
		return core.NewUnauthorizedError("user is not part of the organization")
	}

	role, err := s.RoleRepository.GetRoleByIdentityAndOrganizationIdentity(role_repositories.GetRoleByIdentityAndOrganizationIdentityParams{
		RoleIdentity:         input.RoleIdentity,
		OrganizationIdentity: input.OrganizationIdentity,
	})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if role == nil {
		tx.Rollback()
		return core.NewNotFoundError("role not found")
	}

	if input.Name != nil {
		err = role.ChangeName(*input.Name, &input.UserEditorIdentity)
	}

	if input.Description != nil {
		err = role.ChangeDescription(*input.Description, &input.UserEditorIdentity)
	}

	if input.Permissions != nil {
		role.ClearPermissions(&input.UserEditorIdentity)

		for _, permission := range input.Permissions {
			permission, err := s.PermissionRepository.GetPermissionBySlug(role_repositories.GetPermissionBySlugParams{
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

			role.AddPermission(*permission, &input.UserEditorIdentity)
		}
	}

	err = s.RoleRepository.UpdateRole(role_repositories.UpdateRoleParams{Role: role})
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
