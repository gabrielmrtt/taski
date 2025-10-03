package roleservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
)

type DeleteRoleService struct {
	RoleRepository        rolerepo.RoleRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteRoleService(
	roleRepository rolerepo.RoleRepository,
	transactionRepository core.TransactionRepository,
) *DeleteRoleService {
	return &DeleteRoleService{
		RoleRepository:        roleRepository,
		TransactionRepository: transactionRepository,
	}
}

type DeleteRoleInput struct {
	RoleIdentity         core.Identity
	OrganizationIdentity core.Identity
}

func (i DeleteRoleInput) Validate() error {
	return nil
}

func (s *DeleteRoleService) Execute(input DeleteRoleInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.RoleRepository.SetTransaction(tx)

	role, err := s.RoleRepository.GetRoleByIdentityAndOrganizationIdentity(rolerepo.GetRoleByIdentityAndOrganizationIdentityParams{
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

	role.Delete()

	err = s.RoleRepository.UpdateRole(rolerepo.UpdateRoleParams{Role: role})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	err = s.RoleRepository.ChangeRoleUsersToDefault(rolerepo.ChangeRoleUsersToDefaultParams{
		RoleIdentity:    input.RoleIdentity,
		DefaultRoleSlug: "default",
	})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	return nil
}
