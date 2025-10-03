package userservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user "github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type RecoverUserPasswordService struct {
	UserRepository             userrepo.UserRepository
	PasswordRecoveryRepository userrepo.PasswordRecoveryRepository
	TransactionRepository      core.TransactionRepository
}

func NewRecoverUserPasswordService(
	userRepository userrepo.UserRepository,
	passwordRecoveryRepository userrepo.PasswordRecoveryRepository,
	transactionRepository core.TransactionRepository,
) *RecoverUserPasswordService {
	return &RecoverUserPasswordService{
		UserRepository:             userRepository,
		PasswordRecoveryRepository: passwordRecoveryRepository,
		TransactionRepository:      transactionRepository,
	}
}

type RecoverUserPasswordInput struct {
	PasswordRecoveryToken string
	Password              string
	PasswordConfirmation  string
}

func (i RecoverUserPasswordInput) Validate() error {
	var fields []core.InvalidInputErrorField

	_, err := user.NewPassword(i.Password)
	if err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "password",
			Error: err.Error(),
		})
	}

	if i.Password != i.PasswordConfirmation {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "password_confirmation",
			Error: "password confirmation does not match",
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *RecoverUserPasswordService) Execute(input RecoverUserPasswordInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)
	s.PasswordRecoveryRepository.SetTransaction(tx)

	passwordRecovery, err := s.PasswordRecoveryRepository.GetPasswordRecoveryByToken(userrepo.GetPasswordRecoveryByTokenParams{Token: input.PasswordRecoveryToken})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if passwordRecovery == nil {
		tx.Rollback()
		return core.NewNotFoundError("password recovery not found")
	}

	if passwordRecovery.IsUsed() {
		tx.Rollback()
		return core.NewAlreadyExistsError("password recovery already used")
	}

	if passwordRecovery.IsExpired() {
		tx.Rollback()
		return core.NewAlreadyExistsError("password recovery expired")
	}

	usr, err := s.UserRepository.GetUserByIdentity(userrepo.GetUserByIdentityParams{UserIdentity: passwordRecovery.UserIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if usr == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	err = usr.ChangeCredentialsPassword(input.Password)
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	err = s.UserRepository.UpdateUser(userrepo.UpdateUserParams{User: usr})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	passwordRecovery.Use()

	err = s.PasswordRecoveryRepository.UpdatePasswordRecovery(userrepo.UpdatePasswordRecoveryParams{PasswordRecovery: passwordRecovery})
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
