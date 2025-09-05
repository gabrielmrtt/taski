package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type UpdateUserCredentialsService struct {
	UserRepository        user_core.UserRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateUserCredentialsService(
	userRepository user_core.UserRepository,
	transactionRepository core.TransactionRepository,
) *UpdateUserCredentialsService {
	return &UpdateUserCredentialsService{
		UserRepository:        userRepository,
		TransactionRepository: transactionRepository,
	}
}

type UpdateUserCredentialsInput struct {
	Name        *string
	Email       *string
	PhoneNumber *string
}

func (s *UpdateUserCredentialsService) Execute(userIdentity core.Identity, input UpdateUserCredentialsInput) error {
	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)

	user, err := s.UserRepository.GetUserByIdentity(user_core.GetUserByIdentityParams{
		Identity: userIdentity,
		Include: map[string]any{
			"credentials": true,
		},
	})

	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if user == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	if input.Name != nil {
		user.ChangeCredentialsName(*input.Name)
	}

	if input.Email != nil {
		err = user.ChangeCredentialsEmail(*input.Email)

		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
	}

	if input.PhoneNumber != nil {
		err = user.ChangeCredentialsPhoneNumber(*input.PhoneNumber)

		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
	}

	err = s.UserRepository.UpdateUser(user)

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
