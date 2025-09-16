package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
)

type ListMyOrganizationInvitesService struct {
	OrganizationRepository organization_repositories.OrganizationRepository
}

func NewListMyOrganizationInvitesService(organizationRepository organization_repositories.OrganizationRepository) *ListMyOrganizationInvitesService {
	return &ListMyOrganizationInvitesService{
		OrganizationRepository: organizationRepository,
	}
}

type ListMyOrganizationInvitesInput struct {
	AuthenticatedUserIdentity core.Identity
	Pagination                *core.PaginationInput
	SortInput                 *core.SortInput
	RelationsInput            core.RelationsInput
}

func (i ListMyOrganizationInvitesInput) Validate() error {
	return nil
}

func (s *ListMyOrganizationInvitesService) Execute(input ListMyOrganizationInvitesInput) (*core.PaginationOutput[organization_core.OrganizationDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizations, err := s.OrganizationRepository.PaginateInvitedOrganizationsBy(organization_repositories.PaginateInvitedOrganizationsParams{
		LoggedUserIdentity: input.AuthenticatedUserIdentity,
		SortInput:          input.SortInput,
		Pagination:         input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var organizationsDto []organization_core.OrganizationDto = make([]organization_core.OrganizationDto, 0)
	for _, organization := range organizations.Data {
		organizationsDto = append(organizationsDto, *organization_core.OrganizationToDto(&organization))
	}

	return &core.PaginationOutput[organization_core.OrganizationDto]{
		Data:    organizationsDto,
		Page:    organizations.Page,
		HasMore: organizations.HasMore,
		Total:   organizations.Total,
	}, nil
}
