package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type VerifyUserRegistrationService struct {
	UserRegistrationRepository user_core.UserRegistrationRepository
	UserRepository             user_core.UserRepository
	TransactionRepository      core.TransactionRepository
}

func NewVerifyUserRegistrationService(
	userRegistrationRepository user_core.UserRegistrationRepository,
	userRepository user_core.UserRepository,
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

func (s *VerifyUserRegistrationService) Execute(input VerifyUserRegistrationInput) error {
	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRegistrationRepository.SetTransaction(tx)

	userRegistration, err := s.UserRegistrationRepository.GetUserRegistrationByToken(user_core.GetUserRegistrationByTokenParams{
		Token:   input.Token,
		Include: map[string]any{},
	})

	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if userRegistration == nil {
		tx.Rollback()
		return core.NewNotFoundError("user registration not found")
	}

	now := datetimeutils.EpochNow()

	if userRegistration.Status == user_core.UserRegistrationStatusExpired || userRegistration.ExpiresAt < now {
		tx.Rollback()
		return core.NewAlreadyExistsError("user registration expired")
	}

	user, err := s.UserRepository.GetUserByIdentity(user_core.GetUserByIdentityParams{
		Identity: userRegistration.UserIdentity,
		Include:  map[string]any{},
	})

	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if user == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	user.Status = user_core.UserStatusActive

	_, err = s.UserRepository.StoreUser(user)

	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	userRegistration.Verify()

	err = s.UserRegistrationRepository.UpdateUserRegistration(userRegistration)

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
