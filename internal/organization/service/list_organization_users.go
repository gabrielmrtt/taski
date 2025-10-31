package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
)

type ListOrganizationUsersService struct {
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
}

func NewListOrganizationUsersService(
	organizationUserRepository organizationrepo.OrganizationUserRepository,
) *ListOrganizationUsersService {
	return &ListOrganizationUsersService{
		OrganizationUserRepository: organizationUserRepository,
	}
}

type ListOrganizationUsersInput struct {
	Filters        organizationrepo.OrganizationUserFilters
	Pagination     core.PaginationInput
	SortInput      core.SortInput
	RelationsInput core.RelationsInput
}

func (i ListOrganizationUsersInput) Validate() error {
	return nil
}

func (s *ListOrganizationUsersService) Execute(input ListOrganizationUsersInput) (*core.PaginationOutput[organization.OrganizationUserDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizationUsers, err := s.OrganizationUserRepository.PaginateOrganizationUsersBy(organizationrepo.PaginateOrganizationUsersParams{
		Filters:        input.Filters,
		SortInput:      input.SortInput,
		Pagination:     input.Pagination,
		RelationsInput: input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	var organizationUsersDto []organization.OrganizationUserDto = make([]organization.OrganizationUserDto, 0)
	for _, organizationUser := range organizationUsers.Data {
		organizationUsersDto = append(organizationUsersDto, *organization.OrganizationUserToDto(&organizationUser))
	}

	return &core.PaginationOutput[organization.OrganizationUserDto]{
		Data:    organizationUsersDto,
		Page:    organizationUsers.Page,
		HasMore: organizationUsers.HasMore,
		Total:   organizationUsers.Total,
	}, nil
}
