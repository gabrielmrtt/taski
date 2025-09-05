package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type DeleteUserService struct {
	UserRepository        user_core.UserRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteUserService(
	userRepository user_core.UserRepository,
	transactionRepository core.TransactionRepository,
) *DeleteUserService {
	return &DeleteUserService{
		UserRepository:        userRepository,
		TransactionRepository: transactionRepository,
	}
}

func (s *DeleteUserService) Execute(userIdentity core.Identity) error {
	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)

	user, err := s.UserRepository.GetUserByIdentity(user_core.GetUserByIdentityParams{
		Identity: userIdentity,
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

	user.Delete()

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
