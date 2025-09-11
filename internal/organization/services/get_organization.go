package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
)

type GetOrganizationService struct {
	OrganizationRepository organization_repositories.OrganizationRepository
}

func NewGetOrganizationService(
	organizationRepository organization_repositories.OrganizationRepository,
) *GetOrganizationService {
	return &GetOrganizationService{
		OrganizationRepository: organizationRepository,
	}
}

type GetOrganizationInput struct {
	OrganizationIdentity core.Identity
}

func (i GetOrganizationInput) Validate() error {
	return nil
}

func (s *GetOrganizationService) Execute(input GetOrganizationInput) (*organization_core.OrganizationDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organization, err := s.OrganizationRepository.GetOrganizationByIdentity(organization_repositories.GetOrganizationByIdentityParams{OrganizationIdentity: input.OrganizationIdentity})
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	if organization == nil {
		return nil, core.NewNotFoundError("organization not found")
	}

	return organization_core.OrganizationToDto(organization), nil
}
