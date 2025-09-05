package user_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type UpdateUserDataService struct {
	UserRepository        user_core.UserRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateUserDataService(
	userRepository user_core.UserRepository,
	transactionRepository core.TransactionRepository,
) *UpdateUserDataService {
	return &UpdateUserDataService{
		UserRepository:        userRepository,
		TransactionRepository: transactionRepository,
	}
}

type UpdateUserDataInput struct {
	DisplayName    *string
	About          *string
	ProfilePicture *core.FileUploadInput
}

func (s *UpdateUserDataService) Execute(userIdentity core.Identity, input UpdateUserDataInput) error {
	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)

	user, err := s.UserRepository.GetUserByIdentity(user_core.GetUserByIdentityParams{
		Identity: userIdentity,
		Include: map[string]any{
			"data": true,
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

	if input.DisplayName != nil {
		user.ChangeUserDataDisplayName(*input.DisplayName)
	}

	if input.About != nil {
		user.ChangeUserDataAbout(input.About)
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
