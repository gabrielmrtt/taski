package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
)

type ListOrganizationsService struct {
	OrganizationRepository organization_core.OrganizationRepository
}

func NewListOrganizationsService(
	organizationRepository organization_core.OrganizationRepository,
) *ListOrganizationsService {
	return &ListOrganizationsService{
		OrganizationRepository: organizationRepository,
	}
}

type ListOrganizationsInput struct {
	Filters    organization_core.OrganizationFilters
	Pagination *core.PaginationInput
	SortInput  *core.SortInput
	Include    map[string]any
}

func (i ListOrganizationsInput) Validate() error {
	return nil
}

func (s *ListOrganizationsService) Execute(input ListOrganizationsInput) (*core.PaginationOutput[organization_core.OrganizationDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizations, err := s.OrganizationRepository.PaginateOrganizationsBy(organization_core.PaginateOrganizationsParams{
		Filters:    input.Filters,
		Include:    input.Include,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var organizationsDto []organization_core.OrganizationDto
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
