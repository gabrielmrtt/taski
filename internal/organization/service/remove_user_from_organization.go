package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type RemoveUserFromOrganizationService struct {
	OrganizationRepository     organizationrepo.OrganizationRepository
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
	UserRepository             userrepo.UserRepository
	TransactionRepository      core.TransactionRepository
}

func NewRemoveUserFromOrganizationService(
	organizationRepository organizationrepo.OrganizationRepository,
	organizationUserRepository organizationrepo.OrganizationUserRepository,
	userRepository userrepo.UserRepository,
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

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})
	if err != nil {
		return err
	}

	if organizationUser == nil {
		return core.NewNotFoundError("organization user not found")
	}

	err = s.OrganizationUserRepository.DeleteOrganizationUser(organizationrepo.DeleteOrganizationUserParams{
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
