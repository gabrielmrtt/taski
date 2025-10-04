package userservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user "github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type GetMeService struct {
	UserRepository userrepo.UserRepository
}

func NewGetMeService(
	userRepository userrepo.UserRepository,
) *GetMeService {
	return &GetMeService{
		UserRepository: userRepository,
	}
}

type GetMeInput struct {
	AuthenticatedUserIdentity core.Identity
	RelationsInput            core.RelationsInput
}

func (i GetMeInput) Validate() error {
	return nil
}

func (s *GetMeService) Execute(input GetMeInput) (*user.UserDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	usr, err := s.UserRepository.GetUserByIdentity(userrepo.GetUserByIdentityParams{UserIdentity: input.AuthenticatedUserIdentity})
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	if usr == nil {
		return nil, core.NewNotFoundError("user not found")
	}

	return user.UserToDto(usr), nil
}
