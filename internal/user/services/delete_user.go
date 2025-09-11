package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
)

type DeleteUserService struct {
	UserRepository        user_repositories.UserRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteUserService(
	userRepository user_repositories.UserRepository,
	transactionRepository core.TransactionRepository,
) *DeleteUserService {
	return &DeleteUserService{
		UserRepository:        userRepository,
		TransactionRepository: transactionRepository,
	}
}

type DeleteUserInput struct {
	UserIdentity core.Identity
}

func (i DeleteUserInput) Validate() error {
	return nil
}

func (s *DeleteUserService) Execute(input DeleteUserInput) error {
	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)

	user, err := s.UserRepository.GetUserByIdentity(user_repositories.GetUserByIdentityParams{UserIdentity: input.UserIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if user == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	user.Delete()

	err = s.UserRepository.UpdateUser(user_repositories.UpdateUserParams{User: user})
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
