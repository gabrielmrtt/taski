package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type RefuseOrganizationUserInvitationService struct {
	OrganizationRepository     organization_core.OrganizationRepository
	OrganizationUserRepository organization_core.OrganizationUserRepository
	UserRepository             user_core.UserRepository
	TransactionRepository      core.TransactionRepository
}

func NewRefuseOrganizationUserInvitationService(
	organizationRepository organization_core.OrganizationRepository,
	organizationUserRepository organization_core.OrganizationUserRepository,
	userRepository user_core.UserRepository,
	transactionRepository core.TransactionRepository,
) *RefuseOrganizationUserInvitationService {
	return &RefuseOrganizationUserInvitationService{
		OrganizationRepository:     organizationRepository,
		OrganizationUserRepository: organizationUserRepository,
		UserRepository:             userRepository,
		TransactionRepository:      transactionRepository,
	}
}

type RefuseOrganizationUserInvitationInput struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

func (i RefuseOrganizationUserInvitationInput) Validate() error {
	return nil
}

func (s *RefuseOrganizationUserInvitationService) Execute(input RefuseOrganizationUserInvitationInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.OrganizationRepository.SetTransaction(tx)
	s.UserRepository.SetTransaction(tx)

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organization_core.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})
	if err != nil {
		return err
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	organizationUser.RefuseInvitation()

	err = s.OrganizationUserRepository.UpdateOrganizationUser(organization_core.UpdateOrganizationUserParams{OrganizationUser: organizationUser})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
