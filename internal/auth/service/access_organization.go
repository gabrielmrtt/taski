package authservice

import (
	"github.com/gabrielmrtt/taski/internal/auth"
	"github.com/gabrielmrtt/taski/internal/core"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
)

type AccessOrganizationService struct {
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
	TokenService               auth.TokenService
}

func NewAccessOrganizationService(organizationUserRepository organizationrepo.OrganizationUserRepository, tokenService auth.TokenService) *AccessOrganizationService {
	return &AccessOrganizationService{
		OrganizationUserRepository: organizationUserRepository,
		TokenService:               tokenService,
	}
}

type AccessOrganizationInput struct {
	AuthenticatedUserIdentity core.Identity
	OrganizationIdentity      core.Identity
}

func (i AccessOrganizationInput) Validate() error {
	return nil
}

func (s *AccessOrganizationService) Execute(input AccessOrganizationInput) (*string, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizationUser, err := s.OrganizationUserRepository.GetOrganizationUserByIdentity(organizationrepo.GetOrganizationUserByIdentityParams{
		OrganizationIdentity: input.OrganizationIdentity,
		UserIdentity:         input.AuthenticatedUserIdentity,
	})
	if err != nil {
		return nil, err
	}

	if organizationUser == nil {
		return nil, core.NewNotFoundError("organization not found")
	}

	organizationUser.Access()

	err = s.OrganizationUserRepository.UpdateOrganizationUser(organizationrepo.UpdateOrganizationUserParams{
		OrganizationUser: organizationUser,
	})
	if err != nil {
		return nil, err
	}

	token, err := s.TokenService.GenerateToken(auth.TokenClaims{
		AuthenticatedUserId:             organizationUser.User.Identity.Public,
		AuthenticatedUserOrganizationId: &organizationUser.OrganizationIdentity.Public,
	})
	if err != nil {
		return nil, err
	}

	return &token, nil
}
