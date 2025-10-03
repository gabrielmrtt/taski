package userservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user "github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type UpdateUserCredentialsService struct {
	UserRepository        userrepo.UserRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateUserCredentialsService(
	userRepository userrepo.UserRepository,
	transactionRepository core.TransactionRepository,
) *UpdateUserCredentialsService {
	return &UpdateUserCredentialsService{
		UserRepository:        userRepository,
		TransactionRepository: transactionRepository,
	}
}

type UpdateUserCredentialsInput struct {
	UserIdentity core.Identity
	Name         *string
	Email        *string
	PhoneNumber  *string
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
		_, err := user.NewEmail(*i.Email)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "email",
				Error: err.Error(),
			})
		}
	}

	if i.PhoneNumber != nil {
		_, err := user.NewPhoneNumber(*i.PhoneNumber)
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

func (s *UpdateUserCredentialsService) Execute(input UpdateUserCredentialsInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)

	usr, err := s.UserRepository.GetUserByIdentity(userrepo.GetUserByIdentityParams{UserIdentity: input.UserIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if usr == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	if input.Name != nil {
		err = usr.ChangeCredentialsName(*input.Name)
		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
	}

	if input.Email != nil {
		err = usr.ChangeCredentialsEmail(*input.Email)
		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
	}

	if input.PhoneNumber != nil {
		err = usr.ChangeCredentialsPhoneNumber(*input.PhoneNumber)
		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
	}

	err = s.UserRepository.UpdateUser(userrepo.UpdateUserParams{User: usr})
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
