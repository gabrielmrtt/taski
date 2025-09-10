package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type GetMeService struct {
	UserRepository user_core.UserRepository
}

func NewGetMeService(
	userRepository user_core.UserRepository,
) *GetMeService {
	return &GetMeService{
		UserRepository: userRepository,
	}
}

type GetMeInput struct {
	LoggedUserIdentity core.Identity
}

func (s *GetMeService) Execute(input GetMeInput) (*user_core.UserDto, error) {
	user, err := s.UserRepository.GetUserByIdentity(user_core.GetUserByIdentityParams{UserIdentity: input.LoggedUserIdentity})
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	if user == nil {
		return nil, core.NewNotFoundError("user not found")
	}

	return user_core.UserToDto(user), nil
}
