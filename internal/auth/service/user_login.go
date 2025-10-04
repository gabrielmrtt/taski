package authservice

import (
	"github.com/gabrielmrtt/taski/internal/auth"
	"github.com/gabrielmrtt/taski/internal/core"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	user "github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type UserLoginService struct {
	UserRepository             userrepo.UserRepository
	OrganizationUserRepository organizationrepo.OrganizationUserRepository
	TokenService               auth.TokenService
}

func NewUserLoginService(
	userRepository userrepo.UserRepository,
	organizationUserRepository organizationrepo.OrganizationUserRepository,
	tokenService auth.TokenService,
) *UserLoginService {
	return &UserLoginService{
		UserRepository:             userRepository,
		OrganizationUserRepository: organizationUserRepository,
		TokenService:               tokenService,
	}
}

type UserLoginInput struct {
	Email    string
	Password string
}

func (i UserLoginInput) Validate() error {
	var fields []core.InvalidInputErrorField

	_, err := user.NewEmail(i.Email)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "email",
			Error: err.Error(),
		})
	}

	if i.Password == "" {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "password",
			Error: "password is required",
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UserLoginService) Execute(input UserLoginInput) (*auth.UserAuthDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	usr, err := s.UserRepository.GetUserByEmail(userrepo.GetUserByEmailParams{
		Email: input.Email,
	})
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	if usr == nil {
		return nil, core.NewNotFoundError("user not found")
	}

	if usr.IsDeleted() {
		return nil, core.NewNotFoundError("user not found")
	}

	if !usr.CheckPassword(input.Password) {
		return nil, core.NewUnauthorizedError("invalid password")
	}

	if usr.IsInactive() || usr.IsUnverified() {
		return nil, core.NewUnauthorizedError("user is not activated")
	}

	organizationUser, err := s.OrganizationUserRepository.GetLastAccessedOrganizationUserByUserIdentity(organizationrepo.GetLastAccessedOrganizationUserByUserIdentityParams{
		UserIdentity: usr.Identity,
	})
	if err != nil {
		return nil, err
	}

	var organizationIdentity *core.Identity = nil
	if organizationUser != nil {
		organizationIdentity = &organizationUser.OrganizationIdentity
	}

	var organizationIdentityPublic *string = nil
	if organizationIdentity != nil {
		organizationIdentityPublic = &organizationIdentity.Public
	}

	jwtToken, err := s.TokenService.GenerateToken(auth.TokenClaims{
		AuthenticatedUserId:             usr.Identity.Public,
		AuthenticatedUserOrganizationId: organizationIdentityPublic,
	})
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	return auth.UserAuthToDto(usr, jwtToken, organizationIdentityPublic), nil
}
