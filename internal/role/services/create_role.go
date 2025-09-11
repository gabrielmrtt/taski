package role_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
)

type CreateRoleService struct {
	RoleRepository        role_repositories.RoleRepository
	PermissionRepository  role_repositories.PermissionRepository
	TransactionRepository core.TransactionRepository
}

func NewCreateRoleService(
	roleRepository role_repositories.RoleRepository,
	permissionRepository role_repositories.PermissionRepository,
	transactionRepository core.TransactionRepository,
) *CreateRoleService {
	return &CreateRoleService{
		RoleRepository:        roleRepository,
		PermissionRepository:  permissionRepository,
		TransactionRepository: transactionRepository,
	}
}

type CreateRoleInput struct {
	Name                 string
	Description          string
	Permissions          []string
	OrganizationIdentity core.Identity
	UserCreatorIdentity  core.Identity
}

func (i CreateRoleInput) Validate() error {
	var fields []core.InvalidInputErrorField

	_, err := core.NewName(i.Name)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "name",
			Error: err.Error(),
		})
	}

	_, err = core.NewDescription(i.Description)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "description",
			Error: err.Error(),
		})
	}

	duplicates := make(map[string]struct{})

	for _, permission := range i.Permissions {
		if _, ok := duplicates[permission]; ok {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "permissions",
				Error: "duplicate permission",
			})
		}
		duplicates[permission] = struct{}{}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *CreateRoleService) Execute(input CreateRoleInput) (*role_core.RoleDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.RoleRepository.SetTransaction(tx)
	s.PermissionRepository.SetTransaction(tx)

	organizationHasUser, err := s.RoleRepository.CheckIfOrganizatonHasUser(input.OrganizationIdentity, input.UserCreatorIdentity)
	if err != nil {
		return nil, err
	}

	if !organizationHasUser {
		return nil, core.NewUnauthorizedError("user is not part of the organization")
	}

	permissions := make([]role_core.Permission, 0)
	for _, permissionSlug := range input.Permissions {
		permission, err := s.PermissionRepository.GetPermissionBySlug(role_repositories.GetPermissionBySlugParams{
			Slug: permissionSlug,
		})
		if err != nil {
			return nil, err
		}
		if permission == nil {
			return nil, core.NewNotFoundError("permission not found")
		}
		permissions = append(permissions, *permission)
	}

	role, err := role_core.NewRole(role_core.NewRoleInput{
		Name:                 input.Name,
		Description:          input.Description,
		Permissions:          permissions,
		OrganizationIdentity: &input.OrganizationIdentity,
		UserCreatorIdentity:  &input.UserCreatorIdentity,
		IsSystemDefault:      false,
	})
	if err != nil {
		return nil, err
	}

	role, err = s.RoleRepository.StoreRole(role_repositories.StoreRoleParams{Role: role})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return role_core.RoleToDto(role), nil
}
