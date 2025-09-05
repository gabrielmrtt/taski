package user_services

import (
	"time"

	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type RegisterUserService struct {
	UserRepository             user_core.UserRepository
	UserRegistrationRepository user_core.UserRegistrationRepository
	TransactionRepository      core.TransactionRepository
}

func NewRegisterUserService(
	userRepository user_core.UserRepository,
	userRegistrationRepository user_core.UserRegistrationRepository,
	transactionRepository core.TransactionRepository,
) *RegisterUserService {
	return &RegisterUserService{
		UserRepository:             userRepository,
		UserRegistrationRepository: userRegistrationRepository,
		TransactionRepository:      transactionRepository,
	}
}

type RegisterUserInput struct {
	Name        string
	Email       string
	Password    string
	PhoneNumber *string
}

func (i RegisterUserInput) Validate() error {
	var fields []core.InvalidInputErrorField

	_, err := core.NewName(i.Name)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "name",
			Error: err.Error(),
		})
	}

	_, err = user_core.NewEmail(i.Email)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "email",
			Error: err.Error(),
		})
	}

	_, err = user_core.NewPassword(i.Password)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "password",
			Error: err.Error(),
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *RegisterUserService) Execute(input RegisterUserInput) (*user_core.UserDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)
	s.UserRegistrationRepository.SetTransaction(tx)

	userAlreadyExists, err := s.UserRepository.GetUserByEmail(user_core.GetUserByEmailParams{
		Email: input.Email,
		Include: map[string]any{
			"credentials": true,
		},
	})

	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	if userAlreadyExists != nil {
		tx.Rollback()
		return nil, core.NewAlreadyExistsError("user with this email already exists")
	}

	user, err := user_core.NewUser(user_core.NewUserInput{
		Name:        input.Name,
		Email:       input.Email,
		Password:    input.Password,
		PhoneNumber: input.PhoneNumber,
	})

	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	_, err = s.UserRepository.StoreUser(user)

	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	userRegistration, err := user_core.NewUserRegistration(user.Identity, 48*time.Hour)

	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	_, err = s.UserRegistrationRepository.StoreUserRegistration(userRegistration)

	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	return user_core.UserToDto(user), nil
}
