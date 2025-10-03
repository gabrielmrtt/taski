package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
)

type DeleteOrganizationService struct {
	OrganizationRepository organizationrepo.OrganizationRepository
	TransactionRepository  core.TransactionRepository
}

func NewDeleteOrganizationService(
	organizationRepository organizationrepo.OrganizationRepository,
	transactionRepository core.TransactionRepository,
) *DeleteOrganizationService {
	return &DeleteOrganizationService{
		OrganizationRepository: organizationRepository,
		TransactionRepository:  transactionRepository,
	}
}

type DeleteOrganizationInput struct {
	OrganizationIdentity core.Identity
}

func (i DeleteOrganizationInput) Validate() error {
	return nil
}

func (s *DeleteOrganizationService) Execute(input DeleteOrganizationInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.OrganizationRepository.SetTransaction(tx)

	organization, err := s.OrganizationRepository.GetOrganizationByIdentity(organizationrepo.GetOrganizationByIdentityParams{OrganizationIdentity: input.OrganizationIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if organization == nil {
		tx.Rollback()
		return core.NewNotFoundError("organization not found")
	}

	organization.Delete()

	err = s.OrganizationRepository.UpdateOrganization(organizationrepo.UpdateOrganizationParams{Organization: organization})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	return nil
}
