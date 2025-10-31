package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
)

type ListMyOrganizationInvitesService struct {
	OrganizationRepository organizationrepo.OrganizationRepository
}

func NewListMyOrganizationInvitesService(organizationRepository organizationrepo.OrganizationRepository) *ListMyOrganizationInvitesService {
	return &ListMyOrganizationInvitesService{
		OrganizationRepository: organizationRepository,
	}
}

type ListMyOrganizationInvitesInput struct {
	AuthenticatedUserIdentity core.Identity
	Pagination                core.PaginationInput
	SortInput                 core.SortInput
	RelationsInput            core.RelationsInput
}

func (i ListMyOrganizationInvitesInput) Validate() error {
	return nil
}

func (s *ListMyOrganizationInvitesService) Execute(input ListMyOrganizationInvitesInput) (*core.PaginationOutput[organization.OrganizationDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizations, err := s.OrganizationRepository.PaginateInvitedOrganizationsBy(organizationrepo.PaginateInvitedOrganizationsParams{
		AuthenticatedUserIdentity: input.AuthenticatedUserIdentity,
		SortInput:                 input.SortInput,
		Pagination:                input.Pagination,
		RelationsInput:            input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	var organizationsDto []organization.OrganizationDto = make([]organization.OrganizationDto, 0)
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
