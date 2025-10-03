package userservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type VerifyUserRegistrationService struct {
	UserRegistrationRepository userrepo.UserRegistrationRepository
	UserRepository             userrepo.UserRepository
	TransactionRepository      core.TransactionRepository
}

func NewVerifyUserRegistrationService(
	userRegistrationRepository userrepo.UserRegistrationRepository,
	userRepository userrepo.UserRepository,
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

	userRegistration, err := s.UserRegistrationRepository.GetUserRegistrationByToken(userrepo.GetUserRegistrationByTokenParams{Token: input.Token})
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

	usr, err := s.UserRepository.GetUserByIdentity(userrepo.GetUserByIdentityParams{UserIdentity: userRegistration.UserIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if usr == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	usr.Status = user.UserStatusActive

	err = s.UserRepository.UpdateUser(userrepo.UpdateUserParams{User: usr})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	userRegistration.Verify()

	err = s.UserRegistrationRepository.UpdateUserRegistration(userrepo.UpdateUserRegistrationParams{UserRegistration: userRegistration})
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
