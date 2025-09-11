package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
)

type GetMeService struct {
	UserRepository user_repositories.UserRepository
}

func NewGetMeService(
	userRepository user_repositories.UserRepository,
) *GetMeService {
	return &GetMeService{
		UserRepository: userRepository,
	}
}

type GetMeInput struct {
	LoggedUserIdentity core.Identity
}

func (s *GetMeService) Execute(input GetMeInput) (*user_core.UserDto, error) {
	user, err := s.UserRepository.GetUserByIdentity(user_repositories.GetUserByIdentityParams{UserIdentity: input.LoggedUserIdentity})
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	if user == nil {
		return nil, core.NewNotFoundError("user not found")
	}

	return user_core.UserToDto(user), nil
}
