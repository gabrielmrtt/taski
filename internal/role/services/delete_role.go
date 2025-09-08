package role_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
)

type DeleteRoleService struct {
	RoleRepository        role_core.RoleRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteRoleService(
	roleRepository role_core.RoleRepository,
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
	UserDeleterIdentity  core.Identity
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

	organizationHasUser, err := s.RoleRepository.CheckIfOrganizatonHasUser(input.OrganizationIdentity, input.UserDeleterIdentity)
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if !organizationHasUser {
		tx.Rollback()
		return core.NewUnauthorizedError("user is not part of the organization")
	}

	role, err := s.RoleRepository.GetRoleByIdentityAndOrganizationIdentity(role_core.GetRoleByIdentityAndOrganizationIdentityParams{
		Identity:             input.RoleIdentity,
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

	err = s.RoleRepository.UpdateRole(role)
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	err = s.RoleRepository.ChangeRoleUsersToDefault(input.RoleIdentity, "default")
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	return nil
}
