package userservice

import (
	"time"

	"github.com/gabrielmrtt/taski/internal/core"
	user "github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type RegisterUserService struct {
	UserRepository             userrepo.UserRepository
	UserRegistrationRepository userrepo.UserRegistrationRepository
	TransactionRepository      core.TransactionRepository
}

func NewRegisterUserService(
	userRepository userrepo.UserRepository,
	userRegistrationRepository userrepo.UserRegistrationRepository,
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

	_, err = user.NewEmail(i.Email)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "email",
			Error: err.Error(),
		})
	}

	_, err = user.NewPassword(i.Password)
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

func (s *RegisterUserService) Execute(input RegisterUserInput) (*user.UserDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)
	s.UserRegistrationRepository.SetTransaction(tx)

	userAlreadyExists, err := s.UserRepository.GetUserByEmail(userrepo.GetUserByEmailParams{Email: input.Email})
	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	if userAlreadyExists != nil {
		tx.Rollback()
		return nil, core.NewAlreadyExistsError("user with this email already exists")
	}

	usr, err := user.NewUser(user.NewUserInput{
		Name:        input.Name,
		Email:       input.Email,
		Password:    input.Password,
		PhoneNumber: input.PhoneNumber,
	})
	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	_, err = s.UserRepository.StoreUser(userrepo.StoreUserParams{User: usr})
	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	userRegistration, err := user.NewUserRegistration(usr.Identity, 48*time.Hour)
	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	_, err = s.UserRegistrationRepository.StoreUserRegistration(userrepo.StoreUserRegistrationParams{UserRegistration: userRegistration})
	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, core.NewInternalError(err.Error())
	}

	return user.UserToDto(usr), nil
}
