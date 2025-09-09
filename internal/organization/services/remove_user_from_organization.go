package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type RemoveUserFromOrganizationService struct {
	OrganizationRepository     organization_core.OrganizationRepository
	OrganizationUserRepository organization_core.OrganizationUserRepository
	UserRepository             user_core.UserRepository
	TransactionRepository      core.TransactionRepository
}

func NewRemoveUserFromOrganizationService(
	organizationRepository organization_core.OrganizationRepository,
	organizationUserRepository organization_core.OrganizationUserRepository,
	userRepository user_core.UserRepository,
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

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(input.OrganizationIdentity, input.UserIdentity)

	if err != nil {
		return err
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	err = s.OrganizationUserRepository.DeleteOrganizationUser(input.OrganizationIdentity, input.UserIdentity)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
