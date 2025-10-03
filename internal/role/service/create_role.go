package roleservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/role"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
)

type CreateRoleService struct {
	RoleRepository        rolerepo.RoleRepository
	PermissionRepository  rolerepo.PermissionRepository
	TransactionRepository core.TransactionRepository
}

func NewCreateRoleService(
	roleRepository rolerepo.RoleRepository,
	permissionRepository rolerepo.PermissionRepository,
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

func (s *CreateRoleService) Execute(input CreateRoleInput) (*role.RoleDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.RoleRepository.SetTransaction(tx)
	s.PermissionRepository.SetTransaction(tx)

	permissions := make([]role.Permission, 0)
	for _, permissionSlug := range input.Permissions {
		permission, err := s.PermissionRepository.GetPermissionBySlug(rolerepo.GetPermissionBySlugParams{
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

	rol, err := role.NewRole(role.NewRoleInput{
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

	rol, err = s.RoleRepository.StoreRole(rolerepo.StoreRoleParams{Role: rol})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return role.RoleToDto(rol), nil
}
