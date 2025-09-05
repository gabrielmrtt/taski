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

func (i UpdateUserCredentialsInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Name != nil {
		_, err := core.NewName(*i.Name)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "name",
				Error: err.Error(),
			})
		}
	}

	if i.Email != nil {
		_, err := user_core.NewEmail(*i.Email)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "email",
				Error: err.Error(),
			})
		}
	}

	if i.PhoneNumber != nil {
		_, err := user_core.NewPhoneNumber(*i.PhoneNumber)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "phone_number",
				Error: err.Error(),
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateUserCredentialsService) Execute(userIdentity core.Identity, input UpdateUserCredentialsInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

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
		err = user.ChangeCredentialsName(*input.Name)
		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
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
