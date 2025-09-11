package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
)

type UpdateOrganizationUserService struct {
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
	RoleRepository             role_repositories.RoleRepository
	TransactionRepository      core.TransactionRepository
}

func NewUpdateOrganizationUserService(
	organizationUserRepository organization_repositories.OrganizationUserRepository,
	roleRepository role_repositories.RoleRepository,
	transactionRepository core.TransactionRepository,
) *UpdateOrganizationUserService {
	return &UpdateOrganizationUserService{
		OrganizationUserRepository: organizationUserRepository,
		RoleRepository:             roleRepository,
		TransactionRepository:      transactionRepository,
	}
}

type UpdateOrganizationUserInput struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
	RoleIdentity         *core.Identity
	Status               *organization_core.OrganizationUserStatuses
}

func (i UpdateOrganizationUserInput) Validate() error {
	return nil
}

func (s *UpdateOrganizationUserService) Execute(input UpdateOrganizationUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.OrganizationUserRepository.SetTransaction(tx)
	s.RoleRepository.SetTransaction(tx)

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	if input.RoleIdentity != nil {
		role, err := s.RoleRepository.GetRoleByIdentityAndOrganizationIdentity(role_repositories.GetRoleByIdentityAndOrganizationIdentityParams{
			RoleIdentity:         *input.RoleIdentity,
			OrganizationIdentity: input.OrganizationIdentity,
		})
		if err != nil {
			return core.NewInternalError(err.Error())
		}

		if role == nil {
			return core.NewNotFoundError("role not found")
		}

		organizationUser.ChangeRole(role)
	}

	if input.Status != nil {
		if *input.Status == organization_core.OrganizationUserStatusActive && organizationUser.IsInactive() {
			organizationUser.Activate()
		} else if *input.Status == organization_core.OrganizationUserStatusInactive && organizationUser.IsActive() {
			organizationUser.Deactivate()
		} else {
			return core.NewInvalidInputError("invalid input", []core.InvalidInputErrorField{
				{
					Field: "status",
					Error: "valid statuses are: active, inactive",
				},
			})
		}
	}

	err = s.OrganizationUserRepository.UpdateOrganizationUser(organization_repositories.UpdateOrganizationUserParams{
		OrganizationUser: organizationUser,
	})
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	return nil
}
