package user_services

import (
	"slices"
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	storage_core "github.com/gabrielmrtt/taski/internal/storage"
	storage_services "github.com/gabrielmrtt/taski/internal/storage/services"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type UpdateUserDataService struct {
	UserRepository         user_core.UserRepository
	TransactionRepository  core.TransactionRepository
	UploadedFileRepository storage_core.UploadedFileRepository
	StorageRepository      storage_core.StorageRepository
}

func NewUpdateUserDataService(
	userRepository user_core.UserRepository,
	transactionRepository core.TransactionRepository,
	uploadedFileRepository storage_core.UploadedFileRepository,
	storageRepository storage_core.StorageRepository,
) *UpdateUserDataService {
	return &UpdateUserDataService{
		UserRepository:         userRepository,
		TransactionRepository:  transactionRepository,
		UploadedFileRepository: uploadedFileRepository,
		StorageRepository:      storageRepository,
	}
}

type UpdateUserDataInput struct {
	UserIdentity   core.Identity
	DisplayName    *string
	About          *string
	ProfilePicture *core.FileInput
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
		acceptedMimeTypes := storage_core.GetSupportedImageMimeTypes()

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

func (s *UpdateUserDataService) Execute(input UpdateUserDataInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)

	user, err := s.UserRepository.GetUserByIdentity(user_core.GetUserByIdentityParams{
		Identity: input.UserIdentity,
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

	if input.ProfilePicture != nil {
		uploadedFile, err := storage_services.NewUploadFileService(s.UploadedFileRepository, s.StorageRepository).Execute(storage_services.UploadFileInput{
			File:       *input.ProfilePicture,
			Directory:  "users/" + input.UserIdentity.Internal.String() + "/profile_picture",
			UploadedBy: input.UserIdentity,
		})

		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}

		if uploadedFile == nil {
			tx.Rollback()
			return core.NewInternalError("failed to upload file")
		}

		user.ChangeUserDataProfilePicture(&uploadedFile.Identity)
	} else {
		storage_services.NewDeleteFileByIdentityService(s.UploadedFileRepository, s.StorageRepository).Execute(input.UserIdentity)
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
