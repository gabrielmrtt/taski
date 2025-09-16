package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
)

type UserLoginService struct {
	UserRepository user_repositories.UserRepository
}

func NewUserLoginService(
	userRepository user_repositories.UserRepository,
) *UserLoginService {
	return &UserLoginService{
		UserRepository: userRepository,
	}
}

type UserLoginInput struct {
	Email    string
	Password string
}

func (i UserLoginInput) Validate() error {
	var fields []core.InvalidInputErrorField

	_, err := user_core.NewEmail(i.Email)
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

func (s *UserLoginService) Execute(input UserLoginInput) (*user_core.UserLoginDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := s.UserRepository.GetUserByEmail(user_repositories.GetUserByEmailParams{
		Email: input.Email,
	})
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	if user == nil {
		return nil, core.NewNotFoundError("user not found")
	}

	if user.IsDeleted() {
		return nil, core.NewNotFoundError("user not found")
	}

	if !user.CheckPassword(input.Password) {
		return nil, core.NewUnauthorizedError("invalid password")
	}

	if user.IsInactive() || user.IsUnverified() {
		return nil, core.NewUnauthorizedError("user is not activated")
	}

	jwtToken, err := user_http_middlewares.GenerateJwtToken(user.Identity)
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	return user_core.UserLoginToDto(user, jwtToken), nil
}
