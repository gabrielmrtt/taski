package userservice

import (
	"time"

	"github.com/gabrielmrtt/taski/internal/core"
	user "github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type ForgotUserPasswordService struct {
	UserRepository             userrepo.UserRepository
	PasswordRecoveryRepository userrepo.PasswordRecoveryRepository
	TransactionRepository      core.TransactionRepository
}

func NewForgotUserPasswordService(
	userRepository userrepo.UserRepository,
	passwordRecoveryRepository userrepo.PasswordRecoveryRepository,
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

	_, err := user.NewEmail(i.Email)
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

	usr, err := s.UserRepository.GetUserByEmail(userrepo.GetUserByEmailParams{Email: input.Email})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if usr == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	passwordRecovery, err := user.NewPasswordRecovery(usr.Identity, 48*time.Hour)
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	_, err = s.PasswordRecoveryRepository.StorePasswordRecovery(userrepo.StorePasswordRecoveryParams{PasswordRecovery: passwordRecovery})
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
