package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
)

type VerifyUserRegistrationService struct {
	UserRegistrationRepository user_repositories.UserRegistrationRepository
	UserRepository             user_repositories.UserRepository
	TransactionRepository      core.TransactionRepository
}

func NewVerifyUserRegistrationService(
	userRegistrationRepository user_repositories.UserRegistrationRepository,
	userRepository user_repositories.UserRepository,
	transactionRepository core.TransactionRepository,
) *VerifyUserRegistrationService {
	return &VerifyUserRegistrationService{
		UserRegistrationRepository: userRegistrationRepository,
		UserRepository:             userRepository,
		TransactionRepository:      transactionRepository,
	}
}

type VerifyUserRegistrationInput struct {
	Token string
}

func (i VerifyUserRegistrationInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Token == "" {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "token",
			Error: "token is required",
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *VerifyUserRegistrationService) Execute(input VerifyUserRegistrationInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRegistrationRepository.SetTransaction(tx)

	userRegistration, err := s.UserRegistrationRepository.GetUserRegistrationByToken(user_repositories.GetUserRegistrationByTokenParams{Token: input.Token})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if userRegistration == nil {
		tx.Rollback()
		return core.NewNotFoundError("user registration not found")
	}

	if userRegistration.IsExpired() {
		tx.Rollback()
		return core.NewAlreadyExistsError("user registration expired")
	}

	user, err := s.UserRepository.GetUserByIdentity(user_repositories.GetUserByIdentityParams{UserIdentity: userRegistration.UserIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if user == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	user.Status = user_core.UserStatusActive

	err = s.UserRepository.UpdateUser(user_repositories.UpdateUserParams{User: user})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	userRegistration.Verify()

	err = s.UserRegistrationRepository.UpdateUserRegistration(user_repositories.UpdateUserRegistrationParams{UserRegistration: userRegistration})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	return nil
}
