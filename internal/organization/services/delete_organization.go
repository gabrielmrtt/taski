package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
)

type DeleteOrganizationService struct {
	OrganizationRepository organization_repositories.OrganizationRepository
	TransactionRepository  core.TransactionRepository
}

func NewDeleteOrganizationService(
	organizationRepository organization_repositories.OrganizationRepository,
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

	organization, err := s.OrganizationRepository.GetOrganizationByIdentity(organization_repositories.GetOrganizationByIdentityParams{OrganizationIdentity: input.OrganizationIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if organization == nil {
		tx.Rollback()
		return core.NewNotFoundError("organization not found")
	}

	organization.Delete()

	err = s.OrganizationRepository.UpdateOrganization(organization_repositories.UpdateOrganizationParams{Organization: organization})
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
