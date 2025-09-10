package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type InviteUserToOrganizationService struct {
	OrganizationRepository     organization_core.OrganizationRepository
	OrganizationUserRepository organization_core.OrganizationUserRepository
	UserRepository             user_core.UserRepository
	RoleRepository             role_core.RoleRepository
	TransactionRepository      core.TransactionRepository
}

func NewInviteUserToOrganizationService(
	organizationRepository organization_core.OrganizationRepository,
	organizationUserRepository organization_core.OrganizationUserRepository,
	userRepository user_core.UserRepository,
	roleRepository role_core.RoleRepository,
	transactionRepository core.TransactionRepository,
) *InviteUserToOrganizationService {
	return &InviteUserToOrganizationService{
		OrganizationRepository:     organizationRepository,
		OrganizationUserRepository: organizationUserRepository,
		UserRepository:             userRepository,
		RoleRepository:             roleRepository,
		TransactionRepository:      transactionRepository,
	}
}

type InviteUserToOrganizationInput struct {
	OrganizationIdentity core.Identity
	Email                string
	RoleIdentity         core.Identity
}

func (i InviteUserToOrganizationInput) Validate() error {
	return nil
}

func (s *InviteUserToOrganizationService) Execute(input InviteUserToOrganizationInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.OrganizationRepository.SetTransaction(tx)
	s.UserRepository.SetTransaction(tx)

	organization, err := s.OrganizationRepository.GetOrganizationByIdentity(organization_core.GetOrganizationByIdentityParams{
		Identity: input.OrganizationIdentity,
	})
	if err != nil {
		return err
	}

	if organization == nil {
		return core.NewNotFoundError("organization not found")
	}

	user, err := s.UserRepository.GetUserByEmail(user_core.GetUserByEmailParams{
		Email: input.Email,
	})

	if user == nil {
		return core.NewNotFoundError("user not found")
	}

	role, err := s.RoleRepository.GetRoleByIdentity(role_core.GetRoleByIdentityParams{
		Identity: input.RoleIdentity,
	})
	if err != nil {
		return err
	}

	if role == nil {
		return core.NewNotFoundError("role not found")
	}

	var organizationUser *organization_core.OrganizationUser = nil

	organizationUser, err = s.OrganizationUserRepository.GetOrganizationUserByIdentity(input.OrganizationIdentity, user.Identity)
	if err != nil {
		return err
	}

	if organizationUser == nil {
		organizationUser, err = organization_core.NewOrganizationUser(organization_core.NewOrganizationUserInput{
			OrganizationIdentity: input.OrganizationIdentity,
			User:                 user,
			Role:                 role,
			Status:               organization_core.OrganizationUserStatusInvited,
		})
		if err != nil {
			return err
		}

		organizationUser, err = s.OrganizationUserRepository.CreateOrganizationUser(organizationUser)
		if err != nil {
			return err
		}
	}

	organizationUser.Invite()

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}
