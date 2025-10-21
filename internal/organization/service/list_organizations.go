package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
)

type ListOrganizationsService struct {
	OrganizationRepository organizationrepo.OrganizationRepository
}

func NewListOrganizationsService(
	organizationRepository organizationrepo.OrganizationRepository,
) *ListOrganizationsService {
	return &ListOrganizationsService{
		OrganizationRepository: organizationRepository,
	}
}

type ListOrganizationsInput struct {
	Filters        organizationrepo.OrganizationFilters
	Pagination     core.PaginationInput
	SortInput      core.SortInput
	RelationsInput core.RelationsInput
}

func (i ListOrganizationsInput) Validate() error {
	return nil
}

func (s *ListOrganizationsService) Execute(input ListOrganizationsInput) (*core.PaginationOutput[organization.OrganizationDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizations, err := s.OrganizationRepository.PaginateOrganizationsBy(organizationrepo.PaginateOrganizationsParams{
		Filters:     input.Filters,
		ShowDeleted: false,
		SortInput:   input.SortInput,
		Pagination:  input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var organizationsDto []organization.OrganizationDto = []organization.OrganizationDto{}
	for _, org := range organizations.Data {
		organizationsDto = append(organizationsDto, *organization.OrganizationToDto(&org))
	}

	return &core.PaginationOutput[organization.OrganizationDto]{
		Data:    organizationsDto,
		Page:    organizations.Page,
		HasMore: organizations.HasMore,
		Total:   organizations.Total,
	}, nil
}
