package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type ChangeUserPasswordService struct {
	UserRepository        user_core.UserRepository
	TransactionRepository core.TransactionRepository
}

func NewChangeUserPasswordService(
	userRepository user_core.UserRepository,
	transactionRepository core.TransactionRepository,
) *ChangeUserPasswordService {
	return &ChangeUserPasswordService{
		UserRepository:        userRepository,
		TransactionRepository: transactionRepository,
	}
}

type ChangeUserPasswordInput struct {
	UserIdentity core.Identity
	Password     string
}

func (i ChangeUserPasswordInput) Validate() error {
	var fields []core.InvalidInputErrorField

	_, err := user_core.NewPassword(i.Password)
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

func (s *ChangeUserPasswordService) Execute(input ChangeUserPasswordInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)

	user, err := s.UserRepository.GetUserByIdentity(user_core.GetUserByIdentityParams{UserIdentity: input.UserIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if user == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	err = user.ChangeCredentialsPassword(input.Password)
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	err = s.UserRepository.UpdateUser(user_core.UpdateUserParams{User: user})
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
