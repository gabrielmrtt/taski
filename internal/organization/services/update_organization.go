package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
)

type UpdateOrganizationService struct {
	OrganizationRepository organization_repositories.OrganizationRepository
	TransactionRepository  core.TransactionRepository
}

func NewUpdateOrganizationService(
	organizationRepository organization_repositories.OrganizationRepository,
	transactionRepository core.TransactionRepository,
) *UpdateOrganizationService {
	return &UpdateOrganizationService{
		OrganizationRepository: organizationRepository,
		TransactionRepository:  transactionRepository,
	}
}

type UpdateOrganizationInput struct {
	OrganizationIdentity core.Identity
	Name                 *string
	UserEditorIdentity   core.Identity
}

func (i UpdateOrganizationInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Name != nil {
		_, err := core.NewName(*i.Name)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "name",
				Error: err.Error(),
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateOrganizationService) Execute(input UpdateOrganizationInput) error {
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

	if input.Name != nil {
		err = organization.ChangeName(*input.Name, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
	}

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
