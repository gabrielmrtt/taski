package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
)

type RemoveUserFromOrganizationService struct {
	OrganizationRepository     organization_repositories.OrganizationRepository
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
	UserRepository             user_repositories.UserRepository
	TransactionRepository      core.TransactionRepository
}

func NewRemoveUserFromOrganizationService(
	organizationRepository organization_repositories.OrganizationRepository,
	organizationUserRepository organization_repositories.OrganizationUserRepository,
	userRepository user_repositories.UserRepository,
	transactionRepository core.TransactionRepository,
) *RemoveUserFromOrganizationService {
	return &RemoveUserFromOrganizationService{
		OrganizationRepository:     organizationRepository,
		OrganizationUserRepository: organizationUserRepository,
		UserRepository:             userRepository,
		TransactionRepository:      transactionRepository,
	}
}

type RemoveUserFromOrganizationInput struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

func (i RemoveUserFromOrganizationInput) Validate() error {
	return nil
}

func (s *RemoveUserFromOrganizationService) Execute(input RemoveUserFromOrganizationInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.OrganizationRepository.SetTransaction(tx)
	s.UserRepository.SetTransaction(tx)

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

	err = s.OrganizationUserRepository.DeleteOrganizationUser(organization_repositories.DeleteOrganizationUserParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
