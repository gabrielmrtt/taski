package user_services

import (
	"slices"
	"strings"

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

func (i UpdateUserDataInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.DisplayName != nil {
		_, err := core.NewName(*i.DisplayName)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "display_name",
				Error: err.Error(),
			})
		}
	}

	if i.About != nil {
		_, err := core.NewDescription(*i.About)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "about",
				Error: err.Error(),
			})
		}
	}

	if i.ProfilePicture != nil {
		acceptedMimeTypes := []string{"image/png", "image/jpeg"}

		if !slices.Contains(acceptedMimeTypes, i.ProfilePicture.FileMimeType) {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "profile_picture",
				Error: "invalid file type. supported mime types are: " + strings.Join(acceptedMimeTypes, ", "),
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateUserDataService) Execute(userIdentity core.Identity, input UpdateUserDataInput) error {
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
		err = user.ChangeUserDataDisplayName(*input.DisplayName)
		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
	}

	if input.About != nil {
		err = user.ChangeUserDataAbout(*input.About)
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
