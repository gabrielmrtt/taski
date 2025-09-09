package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type RemoveUserFromOrganizationService struct {
	OrganizationRepository organization_core.OrganizationRepository
	UserRepository         user_core.UserRepository
	TransactionRepository  core.TransactionRepository
}

func NewRemoveUserFromOrganizationService(
	organizationRepository organization_core.OrganizationRepository,
	userRepository user_core.UserRepository,
	transactionRepository core.TransactionRepository,
) *RemoveUserFromOrganizationService {
	return &RemoveUserFromOrganizationService{
		OrganizationRepository: organizationRepository,
		UserRepository:         userRepository,
		TransactionRepository:  transactionRepository,
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

	organizationUser, err := s.OrganizationRepository.GetOrganizationUserByIdentity(input.OrganizationIdentity, input.UserIdentity)

	if err != nil {
		return err
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	err = s.OrganizationRepository.DeleteOrganizationUser(input.OrganizationIdentity, input.UserIdentity)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
