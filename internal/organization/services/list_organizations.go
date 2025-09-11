package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
)

type ListOrganizationsService struct {
	OrganizationRepository organization_repositories.OrganizationRepository
}

func NewListOrganizationsService(
	organizationRepository organization_repositories.OrganizationRepository,
) *ListOrganizationsService {
	return &ListOrganizationsService{
		OrganizationRepository: organizationRepository,
	}
}

type ListOrganizationsInput struct {
	Filters     organization_repositories.OrganizationFilters
	ShowDeleted bool
	Pagination  *core.PaginationInput
	SortInput   *core.SortInput
}

func (i ListOrganizationsInput) Validate() error {
	return nil
}

func (s *ListOrganizationsService) Execute(input ListOrganizationsInput) (*core.PaginationOutput[organization_core.OrganizationDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizations, err := s.OrganizationRepository.PaginateOrganizationsBy(organization_repositories.PaginateOrganizationsParams{
		Filters:     input.Filters,
		ShowDeleted: input.ShowDeleted,
		SortInput:   input.SortInput,
		Pagination:  input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var organizationsDto []organization_core.OrganizationDto = []organization_core.OrganizationDto{}
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
