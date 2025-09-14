package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
)

type AcceptOrganizationUserInvitationService struct {
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
	TransactionRepository      core.TransactionRepository
}

func NewAcceptOrganizationUserInvitationService(
	organizationUserRepository organization_repositories.OrganizationUserRepository,
	transactionRepository core.TransactionRepository,
) *AcceptOrganizationUserInvitationService {
	return &AcceptOrganizationUserInvitationService{
		OrganizationUserRepository: organizationUserRepository,
		TransactionRepository:      transactionRepository,
	}
}

type AcceptOrganizationUserInvitationInput struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

func (i AcceptOrganizationUserInvitationInput) Validate() error {
	return nil
}

func (s *AcceptOrganizationUserInvitationService) Execute(input AcceptOrganizationUserInvitationInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.OrganizationUserRepository.SetTransaction(tx)

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})
	if err != nil {
		return err
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	organizationUser.AcceptInvitation()

	err = s.OrganizationUserRepository.UpdateOrganizationUser(organization_repositories.UpdateOrganizationUserParams{OrganizationUser: organizationUser})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
