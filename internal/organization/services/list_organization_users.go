package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
)

type ListOrganizationUsersService struct {
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
}

func NewListOrganizationUsersService(
	organizationUserRepository organization_repositories.OrganizationUserRepository,
) *ListOrganizationUsersService {
	return &ListOrganizationUsersService{
		OrganizationUserRepository: organizationUserRepository,
	}
}

type ListOrganizationUsersInput struct {
	Filters        organization_repositories.OrganizationUserFilters
	Pagination     *core.PaginationInput
	SortInput      *core.SortInput
	RelationsInput core.RelationsInput
}

func (i ListOrganizationUsersInput) Validate() error {
	return nil
}

func (s *ListOrganizationUsersService) Execute(input ListOrganizationUsersInput) (*core.PaginationOutput[organization_core.OrganizationUserDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizationUsers, err := s.OrganizationUserRepository.PaginateOrganizationUsersBy(organization_repositories.PaginateOrganizationUsersParams{
		Filters:        input.Filters,
		SortInput:      input.SortInput,
		Pagination:     input.Pagination,
		RelationsInput: input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	var organizationUsersDto []organization_core.OrganizationUserDto = make([]organization_core.OrganizationUserDto, 0)
	for _, organizationUser := range organizationUsers.Data {
		organizationUsersDto = append(organizationUsersDto, *organization_core.OrganizationUserToDto(&organizationUser))
	}

	return &core.PaginationOutput[organization_core.OrganizationUserDto]{
		Data:    organizationUsersDto,
		Page:    organizationUsers.Page,
		HasMore: organizationUsers.HasMore,
		Total:   organizationUsers.Total,
	}, nil
}
