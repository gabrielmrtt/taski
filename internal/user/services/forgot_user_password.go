package user_services

import (
	"time"

	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
)

type ForgotUserPasswordService struct {
	UserRepository             user_repositories.UserRepository
	PasswordRecoveryRepository user_repositories.PasswordRecoveryRepository
	TransactionRepository      core.TransactionRepository
}

func NewForgotUserPasswordService(
	userRepository user_repositories.UserRepository,
	passwordRecoveryRepository user_repositories.PasswordRecoveryRepository,
	transactionRepository core.TransactionRepository,
) *ForgotUserPasswordService {
	return &ForgotUserPasswordService{
		UserRepository:             userRepository,
		PasswordRecoveryRepository: passwordRecoveryRepository,
		TransactionRepository:      transactionRepository,
	}
}

type ForgotUserPasswordInput struct {
	Email string
}

func (i ForgotUserPasswordInput) Validate() error {
	var fields []core.InvalidInputErrorField

	_, err := user_core.NewEmail(i.Email)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "email",
			Error: err.Error(),
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *ForgotUserPasswordService) Execute(input ForgotUserPasswordInput) error {
	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)
	s.PasswordRecoveryRepository.SetTransaction(tx)

	user, err := s.UserRepository.GetUserByEmail(user_repositories.GetUserByEmailParams{Email: input.Email})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if user == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	passwordRecovery, err := user_core.NewPasswordRecovery(user.Identity, 48*time.Hour)
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	_, err = s.PasswordRecoveryRepository.StorePasswordRecovery(user_repositories.StorePasswordRecoveryParams{PasswordRecovery: passwordRecovery})
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
