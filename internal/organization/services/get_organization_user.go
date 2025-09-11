package organization_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
)

type GetOrganizationUserService struct {
	OrganizationUserRepository organization_repositories.OrganizationUserRepository
}

func NewGetOrganizationUserService(organizationUserRepository organization_repositories.OrganizationUserRepository) *GetOrganizationUserService {
	return &GetOrganizationUserService{
		OrganizationUserRepository: organizationUserRepository,
	}
}

type GetOrganizationUserInput struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

func (i GetOrganizationUserInput) Validate() error {
	return nil
}

func (s *GetOrganizationUserService) Execute(input GetOrganizationUserInput) (*organization_core.OrganizationUserDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organization_repositories.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})

	if err != nil {
		return nil, err
	}

	if organizationUser == nil {
		return nil, core.NewNotFoundError("organization user not found")
	}

	return organization_core.OrganizationUserToDto(organizationUser), nil
}
