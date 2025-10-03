package organizationservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
)

type GetOrganizationUserService struct {
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
}

func NewGetOrganizationUserService(organizationUserRepository organizationrepo.OrganizationUserRepository) *GetOrganizationUserService {
	return &GetOrganizationUserService{
		OrganizationUserRepository: organizationUserRepository,
	}
}

type GetOrganizationUserInput struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
	RelationsInput       core.RelationsInput
}

func (i GetOrganizationUserInput) Validate() error {
	return nil
}

func (s *GetOrganizationUserService) Execute(input GetOrganizationUserInput) (*organization.OrganizationUserDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.UserIdentity,
	})

	if err != nil {
		return nil, err
	}

	if organizationUser == nil {
		return nil, core.NewNotFoundError("organization user not found")
	}

	return organization.OrganizationUserToDto(organizationUser), nil
}
