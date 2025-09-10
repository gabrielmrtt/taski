package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
)

type ListOrganizationUsersService struct {
	OrganizationUserRepository organization_core.OrganizationUserRepository
}

func NewListOrganizationUsersService(
	organizationUserRepository organization_core.OrganizationUserRepository,
) *ListOrganizationUsersService {
	return &ListOrganizationUsersService{
		OrganizationUserRepository: organizationUserRepository,
	}
}

type ListOrganizationUsersInput struct {
	OrganizationIdentity core.Identity
	Filters              organization_core.OrganizationUserFilters
	Pagination           *core.PaginationInput
	SortInput            *core.SortInput
	Include              map[string]any
}

func (i ListOrganizationUsersInput) Validate() error {
	return nil
}

func (s *ListOrganizationUsersService) Execute(input ListOrganizationUsersInput) (*core.PaginationOutput[organization_core.OrganizationUserDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizationUsers, err := s.OrganizationUserRepository.PaginateOrganizationUsersBy(input.OrganizationIdentity, organization_core.PaginateOrganizationUsersParams{
		Filters:    input.Filters,
		Include:    input.Include,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
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
