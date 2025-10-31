package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
)

type GetOrganizationService struct {
	OrganizationRepository organizationrepo.OrganizationRepository
}

func NewGetOrganizationService(
	organizationRepository organizationrepo.OrganizationRepository,
) *GetOrganizationService {
	return &GetOrganizationService{
		OrganizationRepository: organizationRepository,
	}
}

type GetOrganizationInput struct {
	OrganizationIdentity core.Identity
	RelationsInput       core.RelationsInput
}

func (i GetOrganizationInput) Validate() error {
	return nil
}

func (s *GetOrganizationService) Execute(input GetOrganizationInput) (*organization.OrganizationDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	org, err := s.OrganizationRepository.GetOrganizationByIdentity(organizationrepo.GetOrganizationByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		RelationsInput:       input.RelationsInput,
	})
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	if org == nil {
		return nil, core.NewNotFoundError("organization not found")
	}

	return organization.OrganizationToDto(org), nil
}
