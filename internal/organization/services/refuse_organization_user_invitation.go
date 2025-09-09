package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type RefuseOrganizationUserInvitationService struct {
	OrganizationRepository organization_core.OrganizationRepository
	UserRepository         user_core.UserRepository
	TransactionRepository  core.TransactionRepository
}

func NewRefuseOrganizationUserInvitationService(
	organizationRepository organization_core.OrganizationRepository,
	userRepository user_core.UserRepository,
	transactionRepository core.TransactionRepository,
) *RefuseOrganizationUserInvitationService {
	return &RefuseOrganizationUserInvitationService{
		OrganizationRepository: organizationRepository,
		UserRepository:         userRepository,
		TransactionRepository:  transactionRepository,
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

	organizationUser, err := s.OrganizationRepository.GetOrganizationUserByIdentity(input.OrganizationIdentity, input.UserIdentity)
	if err != nil {
		return err
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	organizationUser.RefuseInvitation()

	err = s.OrganizationRepository.UpdateOrganizationUser(organizationUser)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
