package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
)

type ListMyOrganizationInvitesService struct {
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
}

func NewListMyOrganizationInvitesService(organizationUserRepository organization_repositories.OrganizationUserRepository) *ListMyOrganizationInvitesService {
	return &ListMyOrganizationInvitesService{
		OrganizationUserRepository: organizationUserRepository,
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

func (s *ListMyOrganizationInvitesService) Execute(input ListMyOrganizationInvitesInput) (*core.PaginationOutput[organization_core.OrganizationUserDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	invitedStatus := organization_core.OrganizationUserStatusInvited
	organizationUsers, err := s.OrganizationUserRepository.PaginateOrganizationUsersBy(organization_repositories.PaginateOrganizationUsersParams{
		Filters: organization_repositories.OrganizationUserFilters{
			UserPublicId: &core.ComparableFilter[string]{
				Equals: &input.AuthenticatedUserIdentity.Public,
			},
			Status: &core.ComparableFilter[organization_core.OrganizationUserStatuses]{
				Equals: &invitedStatus,
			},
		},
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
