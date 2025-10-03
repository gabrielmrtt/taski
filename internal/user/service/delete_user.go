package userservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type DeleteUserService struct {
	UserRepository        userrepo.UserRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteUserService(
	userRepository userrepo.UserRepository,
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

	user, err := s.UserRepository.GetUserByIdentity(userrepo.GetUserByIdentityParams{UserIdentity: input.UserIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if user == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	user.Delete()

	err = s.UserRepository.UpdateUser(userrepo.UpdateUserParams{User: user})
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
